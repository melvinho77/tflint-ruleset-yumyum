package rules

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformNamingStandards struct {
	tflint.DefaultRule
}

type terraformNamingStandardsConfig struct {
	Format string `hclext:"format,optional"`
	Custom string `hclext:"custom,optional"`

	CustomFormats map[string]*CustomFormatConfig `hclext:"custom_formats,optional"`

	Variable *BlockFormatConfig `hclext:"variable,block"`
}

// CustomFormatConfig defines a custom format that can be used instead of the predefined formats
type CustomFormatConfig struct {
	Regexp      string `cty:"regex"`
	Description string `cty:"description"`
}

// BlockFormatConfig defines the pre-defined format or custom regular expression to use for each block type.
// (data, module, resource, output, variable, check ,locals)
type BlockFormatConfig struct {
	Format string `hclext:"format,optional"`
	Custom string `hclext:"custom,optional"`
}

// NameValidator contains the regular expression to validate block name, if it was a named format, and the format name / regular expression string
type NameValidator struct {
	Format        string
	IsNamedFormat bool
	Regexp        *regexp.Regexp
}

var predefinedFormats = map[string]*regexp.Regexp{
	"snake_case":       regexp.MustCompile("^[a-z][a-z0-9]*(_[a-z0-9]+)*$"),
	"mixed_snake_case": regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9]*(_[a-zA-Z0-9]+)*$"),
}

// NewTerraformNamingStandards returns a new rule
func NewTerraformNamingStandards() *TerraformNamingStandards {
	return &TerraformNamingStandards{}
}

// Name returns the rule name
func (r *TerraformNamingStandards) Name() string {
	return "terraform_naming_standards"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformNamingStandards) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformNamingStandards) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformNamingStandards) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_naming_standards.md"
}

// Check verifies that top-level attributes in object-type variables follow naming standards.
// This extends the terraform_naming_convention rule from tflint-ruleset-terraform
// (https://github.com/terraform-linters/tflint-ruleset-terraform/blob/v0.11.0/rules/terraform_naming_convention.go).
// Currently, it only checks surface-level fields in object types and do not check on nested attributes.
func (r *TerraformNamingStandards) Check(runner tflint.Runner) error {
	// Defaults to snake_case
	config := &terraformNamingStandardsConfig{}
	config.Format = "snake_case"

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	defaultNameValidator, err := config.getNameValidator()
	if err != nil {
		return fmt.Errorf("invalid default configuration: %v", err)
	}

	var nameValidator *NameValidator

	body, err := runner.GetModuleContent(&hclext.BodySchema{
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

	blocks := body.Blocks.ByType()

	variableBlockName := "variable"
	nameValidator, err = config.Variable.getNameValidator(defaultNameValidator, config, variableBlockName)
	if err != nil {
		return err
	}

	for _, block := range blocks[variableBlockName] {
		// Top-level variable name check
		if err := nameValidator.checkBlock(runner, r, variableBlockName, block.Labels[0], &block.DefRange); err != nil {
			return err
		}

		typeAttr, typeExist := block.Body.Attributes["type"]
		if !typeExist {
			continue
		}

		switch typeExpr := typeAttr.Expr.(type) {
		case *hclsyntax.FunctionCallExpr:
			switch typeExpr.Name {
			case "object":
				if err := checkNestedObjectFields(typeExpr, runner, r, variableBlockName, block.Labels[0], nameValidator, &typeAttr.Range); err != nil {
					return err
				}

			case "map":
				if len(typeExpr.Args) == 1 {
					if nestedObj, ok := typeExpr.Args[0].(*hclsyntax.FunctionCallExpr); ok && nestedObj.Name == "object" {
						if err := checkNestedObjectFields(nestedObj, runner, r, variableBlockName, block.Labels[0], nameValidator, &typeAttr.Range); err != nil {
							return err
						}
					}
				}
			}
		default:
			continue
		}
	}

	return nil
}

// Checks if a block name (resource, data, output...) matches the expected
// naming convention (snake_case, custom_regex),
// if it is invalid, emits an tf linting issue.
func (validator *NameValidator) checkBlock(runner tflint.Runner, r *TerraformNamingStandards, blockTypeName string, blockName string, blockDeclRange *hcl.Range) error {
	if validator == nil {
		return nil
	}

	// Extracts the last node from full path if applicable
	parts := strings.Split(blockName, ".")
	lastNode := parts[len(parts)-1]

	// Validate only the leaf node
	if !validator.Regexp.MatchString(lastNode) {
		var formatType string

		if validator.IsNamedFormat {
			formatType = "format"
		} else {
			formatType = "RegExp"
		}

		return runner.EmitIssue(
			r,
			fmt.Sprintf("%s name `%s` must match the following %s: %s", blockTypeName, blockName, formatType, validator.Format),
			*blockDeclRange,
		)
	}

	return nil
}

func (config *terraformNamingStandardsConfig) getNameValidator() (*NameValidator, error) {
	return getNameValidator(config.Custom, config.Format, config)
}

// Returns NameValidator for each specific block (data, resource, module...)
// If the block has `format` or `custom` attribute specified, overrides the default NameValidator.
func (blockFormatConfig *BlockFormatConfig) getNameValidator(defaultValidator *NameValidator, config *terraformNamingStandardsConfig, blockName string) (*NameValidator, error) {
	validator := defaultValidator

	if blockFormatConfig != nil {
		nameValidator, err := getNameValidator(blockFormatConfig.Custom, blockFormatConfig.Format, config)
		if err != nil {
			return nil, fmt.Errorf("invalid %s configuration: %v", blockName, err)
		}

		validator = nameValidator
	}

	return validator, nil
}

// Builds the NameValidator according to `TerraformNamingStandardsConfig` struct
// 1. If `custom` field is provided, use it directly as the regex.
// 2. Else if `format` is not "none", check in customFormats map first.
// 3. If not found, will check with predefined formats; error if still not found on formats.
func getNameValidator(custom string, format string, config *terraformNamingStandardsConfig) (*NameValidator, error) {
	// Prefers custom format if specified
	if custom != "" {
		return getCustomNameValidator(false, custom, custom)

	} else if format != "none" {
		customFormats := config.CustomFormats
		customFormatConfig, exists := customFormats[format]

		if exists {
			return getCustomNameValidator(true, customFormatConfig.Description, customFormatConfig.Regexp)
		}

		regex, exists := predefinedFormats[strings.ToLower(format)]
		if exists {
			nameValidator := &NameValidator{
				IsNamedFormat: true,
				Format:        format,
				Regexp:        regex,
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
		IsNamedFormat: isNamed,
		Format:        format,
		Regexp:        regex,
	}

	return nameValidator, err
}

// Recursive function that extracts the inner object definition (ObjectConsExpr)
// from complex Terraform type expressions such as:
//   - object({...})
//   - map(object({...}))
//   - list(map(object({...})))
//   - complex variable structures like map(map(object({...})))
func checkNestedObjectFields(expr hclsyntax.Expression, runner tflint.Runner, r *TerraformNamingStandards, blockTypeName, varKey string, nameValidator *NameValidator, defRange *hcl.Range) error {
	objExpr, objExist := unwrapToObjectConsExpr(expr)
	if !objExist {
		return nil
	}

	for _, item := range objExpr.Items {
		var fieldName string

		// Unwrap key expression if it's wrapped
		keyExpr := item.KeyExpr
		if wrappedKey, ok := keyExpr.(*hclsyntax.ObjectConsKeyExpr); ok {
			keyExpr = wrappedKey.Wrapped
		}

		// Extract field name from key expression
		switch key := keyExpr.(type) {
		case *hclsyntax.LiteralValueExpr:
			fieldName = key.Val.AsString()
		case *hclsyntax.ScopeTraversalExpr:
			fieldName = key.Traversal.RootName()

		default:
			continue
		}

		fullPath := fmt.Sprintf("%s.%s", varKey, fieldName)

		if err := nameValidator.checkBlock(runner, r, blockTypeName, fullPath, defRange); err != nil {
			return err
		}

		// Check if the value itself is another object() function call
		if nestedObj, ok := item.ValueExpr.(*hclsyntax.FunctionCallExpr); ok && nestedObj.Name == "object" {
			if err := checkNestedObjectFields(nestedObj, runner, r, blockTypeName, fullPath, nameValidator, defRange); err != nil {
				return err
			}
		}
	}

	return nil
}

// unwrapToObjectConsExpr extracts the inner object({...}) expression, returning the underlying
// ObjectConsExpr which holds the key-value fields.
//
// In Terraform, a type declaration like `type = object({ ... })` is parsed as a FunctionCallExpr
// named "object" with one argument: an ObjectConsExpr containing the key-value fields.
// Thus, we need to extract from it.
func unwrapToObjectConsExpr(expr hclsyntax.Expression) (*hclsyntax.ObjectConsExpr, bool) {
	// Check if the expression is a function call (e.g., object(...))
	functionCall, isFunction := expr.(*hclsyntax.FunctionCallExpr)
	if !isFunction || len(functionCall.Args) != 1 {
		return nil, false
	}

	// Check if the first argument is the object literal { key = value, ... }
	objectExpr, isObject := functionCall.Args[0].(*hclsyntax.ObjectConsExpr)
	if !isObject {
		return nil, false
	}

	return objectExpr, true
}
