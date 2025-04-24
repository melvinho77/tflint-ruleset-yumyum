package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformRequredVariables(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "module with no required variables.",
			Content: `
variable "my_variable" {
  type = string
}
`,
			Config: testTerraformRequiredVariablesConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVariables(),
					Message: "required variable(s) not declared: cloud_creds, module_info, module_tmpl",
					Range: hcl.Range{
						Start: hcl.Pos{Line: 1, Column: 1},
						End:   hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "module with incomplete required variables. (no module_info), and no sensitive attribute on `cloud_creds`.",
			Content: `
variable "cloud_creds" {
  type = string
}

variable "module_tmpl" {
  type = string
}
`,
			Config: testTerraformRequiredVariablesConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVariables(),
					Message: "required variable(s) not declared: module_info",
					Range: hcl.Range{
						Start: hcl.Pos{Line: 1, Column: 1},
						End:   hcl.Pos{Line: 1, Column: 1},
					},
				},
				{
					Rule:    NewTerraformRequiredVariables(),
					Message: "variable `cloud_creds` is missing the `sensitive` attribute",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 23},
					},
				},
			},
		},
		{
			Name: "module with incomplete required variables. (no module_info & module_tmpl), but with sensitive attribute on `cloud_creds`.",
			Content: `
variable "cloud_creds" {
  sensitive = true
  type      = string
}
`,
			Config: testTerraformRequiredVariablesConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVariables(),
					Message: "required variable(s) not declared: module_info, module_tmpl",
					Range: hcl.Range{
						Start: hcl.Pos{Line: 1, Column: 1},
						End:   hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "module with complete required variables and correct attribute.",
			Content: `
variable "cloud_creds" {
  sensitive = true
  type      = string
}

variable "module_tmpl" {
  type      = string
}

variable "module_info" {
  type      = string
}
`,
			Config:   testTerraformRequiredVariablesConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "module with complete required variables, but no senstitive attribute",
			Content: `
variable "cloud_creds" {
  type      = string
}

variable "module_tmpl" {
  type      = string
}

variable "module_info" {
  type      = string
}
`,
			Config: testTerraformRequiredVariablesConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredVariables(),
					Message: "variable `cloud_creds` is missing the `sensitive` attribute",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 23},
					},
				},
			},
		},
	}

	rule := NewTerraformRequiredVariables()
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

const testTerraformRequiredVariablesConfig = `
rule "terraform_required_variables" {
  enabled = true
}
`
