package rules

import (
	"fmt"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformAnyTypeVariables checks whether variables have a type declared
type TerraformAnyTypeVariables struct {
	tflint.DefaultRule
}

type terraformAnyTypeVariablesConfig struct {
	IgnoreVars []string `hclext:"ignore_vars,optional"`
}

// NewTerraformAnyTypeVariables returns a new rule
func NewTerraformAnyTypeVariables() *TerraformAnyTypeVariables {
	return &TerraformAnyTypeVariables{}
}

// Name returns the rule name
func (r *TerraformAnyTypeVariables) Name() string {
	return "terraform_any_type_variables"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformAnyTypeVariables) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformAnyTypeVariables) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformAnyTypeVariables) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_any_type_variables.md"
}

// Check checks whether variables have type
func (r *TerraformAnyTypeVariables) Check(runner tflint.Runner) error {
	config := &terraformAnyTypeVariablesConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	variables, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name: "type",
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, variable := range variables.Blocks {
		// Skip this check if the variable name match any of the string in ignore_vars.
		if slices.Contains(config.IgnoreVars, variable.Labels[0]) {
			continue
		}

		typeAttr, typeExist := variable.Body.Attributes["type"]
		if !typeExist {
			continue
		}

		for _, typeExpr := range typeAttr.Expr.Variables() {
			if typeExpr.RootName() == "any" {
				err := runner.EmitIssue(
					r,
					fmt.Sprintf("variable '%s' has 'any' type declared", variable.Labels[0]),
					typeExpr.SourceRange(),
				)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
