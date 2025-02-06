package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformAnyTypeVariables(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "simple variable without 'any' type",
			Content: `
variable "my_var" {
  type = string
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "complex variable without 'any' type",
			Content: `
variable "my_var" {
  type = object({
    my_key = number
  })
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "ignored variable with 'any' type",
			Content: `
variable "my_var" {
  type = object({
    my_key = number
  })
}

variable "my_ignored_var" {
  type = object({
    my_key = any
  })
}`,
			Config:   testTerraformAnyTypeVariablesConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "simple variable with 'any' type",
			Content: `
variable "my_var" {
  type = any
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformAnyTypeVariables(),
					Message: "variable 'my_var' has 'any' type declared",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 10},
						End:      hcl.Pos{Line: 3, Column: 13},
					},
				},
			},
		},
		{
			Name: "complex variable with 'any' type",
			Content: `
variable "my_var" {
  type = object({
    my_key = any
  })
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformAnyTypeVariables(),
					Message: "variable 'my_var' has 'any' type declared",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 14},
						End:      hcl.Pos{Line: 4, Column: 17},
					},
				},
			},
		},
	}

	rule := NewTerraformAnyTypeVariables()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{
				"main.tf":     test.Content,
				".tflint.hcl": test.Config,
			})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}

const testTerraformAnyTypeVariablesConfig = `
rule "terraform_any_type_variables" {
  enabled     = true

  ignore_vars = ["my_ignored_var"]
}
`
