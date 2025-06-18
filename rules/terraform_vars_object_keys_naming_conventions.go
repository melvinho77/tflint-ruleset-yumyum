package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

type TerraformVarsObjectKeysNamingConventions struct {
	tflint.DefaultRule
}

type terraformVarsObjectKeysNamingConventionsConfig struct {
	Format          string                         `hclext:"format,optional"`
	CustomFormatKey string                         `hclext:"custom_format_key,optional"`
	CustomFormats   map[string]*CustomFormatConfig `hclext:"custom_formats,optional"`
}

// CustomFormatConfig defines a custom format that can be used instead of the predefined formats
type CustomFormatConfig struct {
	Regexp      string `cty:"regex"`
	Description string `cty:"description"`
}

type NameValidator struct {
	Format             string
	IsPredefinedFormat bool
	Regexp             *regexp.Regexp
}

var predefinedFormats = map[string]*regexp.Regexp{
	// snake_case: lowercase letters and digits, separated by underscores.
	// Must start with a lowercase letter.
	// Examples: "example_name", "my_var_1", "foo_bar123"
	"snake_case": regexp.MustCompile("^[a-z][a-z0-9]*(_[a-z0-9]+)*$"),

	// mixed_snake_case: allows both uppercase and lowercase letters,
	// and digits, separated by underscores. Must start with a letter (any case).
	// Examples: "Example_Name", "myVar_2", "My_Example_123"
	"mixed_snake_case": regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*(_[a-zA-Z0-9]+)*$"),
}

// NewTerraformVarsObjectKeysNamingConventions returns a new rule
func NewTerraformVarsObjectKeysNamingConventions() *TerraformVarsObjectKeysNamingConventions {
	return &TerraformVarsObjectKeysNamingConventions{}
}

// Name returns the rule name
func (r *TerraformVarsObjectKeysNamingConventions) Name() string {
	return "terraform_vars_object_keys_naming_conventions"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVarsObjectKeysNamingConventions) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformVarsObjectKeysNamingConventions) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformVarsObjectKeysNamingConventions) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md"
}

// Check verifies that top-level attributes in object-type variables follow naming standards.
// This extends the terraform_naming_convention rule from tflint-ruleset-terraform
// (https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.11.0/rules/terraform_naming_convention.go).
// Currently, it only checks surface-level fields in object types and do not check on nested attributes.
func (r *TerraformVarsObjectKeysNamingConventions) Check(runner tflint.Runner) error {
	// Load rule configuration, defaulting to snake_case
	config := &terraformVarsObjectKeysNamingConventionsConfig{
		Format: "snake_case",
	}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	// Fetch all variable blocks with type attributes
	variables, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "type"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	// Initialize the name validator
	nameValidator, err := config.getNameValidator()
	if err != nil {
		return err
	}

	// Loop through each variable declared
	for _, variable := range variables.Blocks {
		variableName := variable.Labels[0]

		if !nameValidator.Regexp.MatchString(variableName) {
			var formatType string

			if nameValidator.IsPredefinedFormat {
				formatType = "predefined_format"
			} else {
				formatType = "custom_format"
			}

			err := runner.EmitIssue(
				r,
				fmt.Sprintf("variable `%s` must match the following %s: %s", variableName, formatType, nameValidator.Format),
				variable.DefRange,
			)
			if err != nil {
				return err
			}
		}

		typeAttr, ok := variable.Body.Attributes["type"]
		if !ok {
			continue
		}

		// Convert hcl.Expression to hclsyntax.Expression
		syntaxExpr, ok := typeAttr.Expr.(hclsyntax.Expression)
		if !ok {
			continue
		}

		// Recursively validate nested complex types
		if err := checkNestedObjectFields(syntaxExpr, runner, r, variableName, nameValidator, &typeAttr.Range); err != nil {
			return err
		}
	}

	return nil
}

func (nameValidator *NameValidator) validate(
	runner tflint.Runner,
	r *TerraformVarsObjectKeysNamingConventions,
	fullPath string, // Full field path (e.g., "user_info.address.city")
	defRange *hcl.Range, // Range to report the issue in HCL file
) error {
	if nameValidator == nil {
		return nil
	}

	// Extract the first & last node from the full path for validation
	// Example: from "user_info.address.city", extract `user_info` & `city`
	parts := strings.Split(fullPath, ".")
	rootNode := parts[0]
	lastNode := parts[len(parts)-1]

	// Validate the last variable name against the regex or named format
	if !nameValidator.Regexp.MatchString(lastNode) {
		var formatType string

		if nameValidator.IsPredefinedFormat {
			formatType = "predefined_format"
		} else {
			formatType = "custom_format"
		}

		return runner.EmitIssue(
			r,
			fmt.Sprintf(
				"variable `%s` path `%s` - attribute `%s` must match the following %s: %s",
				rootNode,
				fullPath,
				lastNode,
				formatType,
				nameValidator.Format,
			),
			*defRange,
		)
	}

	return nil
}

func (config *terraformVarsObjectKeysNamingConventionsConfig) getNameValidator() (*NameValidator, error) {
	return getNameValidator(config.Format, config.CustomFormatKey, config)
}

// Builds the NameValidator according to `terraformVarsObjectKeysNamingConventionsConfig` struct
// 1. If `format` is not "none", check in customFormats map first.
// 2. If not found, will check with predefined formats (`snake_case`, `mixed_snake_case`); return error if no format is found.
func getNameValidator(format string, customFormatKey string, config *terraformVarsObjectKeysNamingConventionsConfig) (*NameValidator, error) {
	if format != "none" {
		customFormats := config.CustomFormats
		customFormatConfig, exists := customFormats[customFormatKey]

		if exists {
			return getCustomNameValidator(false, customFormatConfig.Description, customFormatConfig.Regexp)
		}

		regex, exists := predefinedFormats[strings.ToLower(format)]

		if exists {
			nameValidator := &NameValidator{
				IsPredefinedFormat: true,
				Format:             format,
				Regexp:             regex,
			}

			return nameValidator, nil
		}

		return nil, fmt.Errorf("`%s` is unsupported format", format)
	}

	return nil, nil
}

// Creates a `NameValidator` struct from `expression` parameter regex string.
func getCustomNameValidator(isNamed bool, format, expression string) (*NameValidator, error) {
	regex, err := regexp.Compile(expression)

	nameValidator := &NameValidator{
		IsPredefinedFormat: isNamed,
		Format:             format,
		Regexp:             regex,
	}

	return nameValidator, err
}

// checkNestedObjectFields recursively validates that all object keys (field names) inside
// Terraform variable type expressions conform to a naming convention.
// It supports deeply nested complex terraform expressions such as:
//   - object({...})
//   - map(object({...}))
//   - list(map(object({...})))
//   - map(map(object({...})))
//   - tuple([object({...})])
func checkNestedObjectFields(
	expr hclsyntax.Expression,
	runner tflint.Runner,
	r *TerraformVarsObjectKeysNamingConventions,
	varKey string,
	nameValidator *NameValidator,
	defRange *hcl.Range,
) error {
	if fnExpr, ok := expr.(*hclsyntax.FunctionCallExpr); ok {
		switch fnExpr.Name {
		case "object":
			// Attempt to unwrap object({...}) structure to extract field definitions
			objExpr, ok := unwrapToObjectConsExpr(fnExpr)
			if !ok {
				return nil
			}

			for _, item := range objExpr.Items {
				// Extract key (field name) from key expression
				fieldName := extractKeyName(item.KeyExpr)
				if fieldName == "" {
					continue
				}

				// Construct full variable path, e.g., "foo.bar.test"
				fullPath := fmt.Sprintf("%s.%s", varKey, fieldName)

				// Check naming convention of the current field name
				if err := nameValidator.validate(runner, r, fullPath, defRange); err != nil {
					return err
				}

				// Recursively check the value expression (in case it's a nested object or complex type)
				if err := checkNestedObjectFields(item.ValueExpr, runner, r, fullPath, nameValidator, defRange); err != nil {
					return err
				}
			}

		case "map", "list", "set", "tuple":
			for _, arg := range fnExpr.Args {
				if err := checkNestedObjectFields(arg, runner, r, varKey, nameValidator, defRange); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// unwrapToObjectConsExpr extracts the underlying ObjectConsExpr from an object() function.
// Terraform represents `type = object({ key = type, ... })` as a FunctionCallExpr with
// one argument: an ObjectConsExpr holding key-value pairs for the object fields.
func unwrapToObjectConsExpr(expr hclsyntax.Expression) (*hclsyntax.ObjectConsExpr, bool) {
	functionCall, isFunction := expr.(*hclsyntax.FunctionCallExpr)
	if !isFunction || functionCall.Name != "object" || len(functionCall.Args) != 1 {
		return nil, false
	}

	// Check if the single argument is an object literal definition: { key = type, ... }
	objExpr, ok := functionCall.Args[0].(*hclsyntax.ObjectConsExpr)

	return objExpr, ok
}

// extractKeyName unwraps and retrieves the field name from an object key expression.
// This supports literal keys (e.g., "name") and traversals (e.g., user.name).
func extractKeyName(keyExpr hclsyntax.Expression) string {
	// Unwrap if it's an ObjectConsKeyExpr (e.g., a quoted key like "field")
	if wrapped, ok := keyExpr.(*hclsyntax.ObjectConsKeyExpr); ok {
		keyExpr = wrapped.Wrapped
	}

	switch key := keyExpr.(type) {
	case *hclsyntax.LiteralValueExpr:
		if key.Val.Type() == cty.String {
			return key.Val.AsString()
		}
	case *hclsyntax.ScopeTraversalExpr:
		return key.Traversal.RootName()
	}

	return ""
}
