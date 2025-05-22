package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformModuleDependencies(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "module with no source.",
			Content: `
module "my_module" {
  name = "my_name"
}
  `,
			Expected: helper.Issues{},
		},
		{
			Name: "local module.",
			Content: `
module "my_module" {
  source = "./test"
  name   = "my_name"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module is not pinned",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test/test-module.git"
  name   = "my_name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleSourceVersion(),
					Message: "module 'my_module' source 'git::https://gitlab.example.com/test/test-module.git' is not pinned (missing ?ref= or ?rev= in the URL).",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 66},
					},
				},
			},
		},
		{
			Name: "git module referenced is default branch.",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test/test-module.git?ref=main"
  name   = "my_name"
}
`,
			Config: testTerraformModuleSourceVersionConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleSourceVersion(),
					Message: "module 'my_module' source 'git::https://gitlab.example.com/test/test-module.git?ref=main' [ref='main'] does not match any allowed_versions pattern",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 75},
					},
				},
			},
		},
		{
			Name: "invalid URL",
			Content: `
module "my_module" {
  source = "git://#{}.com"
  name   = "my_name"
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleSourceVersion(),
					Message: "module 'my_module' source 'git://#{}.com' is not a valid URL",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 27},
					},
				},
			},
		},
		{
			Name: "git module reference is pinned to semver.",
			Content: `
module "my_module" {
  source = "git::gitlab.example.com/test.git?ref=v1.2.3"
  name   = "my_name"
}`,
			Config: `
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "semver"
  default_branches = []
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module reference is pinned to semver (no leading v).",
			Content: `
module "my_module" {
  source = "git://gitlab.example.com/test.git?ref=1.2.3"
  name   = "my_name"
}`,
			Config: `
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "semver"
  default_branches = []
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "git module referece is pinned with valid pattern for bugfix branch.",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test.git?ref=bugfix/1234"
  name   = "my_name"
}`,
			Config:   testTerraformModuleSourceVersionConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "git module referece is pinned with valid pattern for feature branch.",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test.git?ref=feature/1234"
  name   = "my_name"
}`,
			Config:   testTerraformModuleSourceVersionConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "git module referece is pinned with invalid pattern for bugfix branch.",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test.git?ref=bugfix/xxx"
  name   = "my_name"
}`,
			Config: testTerraformModuleSourceVersionConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleSourceVersion(),
					Message: "module 'my_module' source 'git::https://gitlab.example.com/test.git?ref=bugfix/xxx' [ref='bugfix/xxx'] does not match any allowed_versions pattern",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 69},
					},
				},
			},
		},
		{
			Name: "git module referece is pinned with invalid pattern for feature branch.",
			Content: `
module "my_module" {
  source = "git::https://gitlab.example.com/test.git?ref=feature/xxx"
  name   = "my_name"
}`,
			Config: testTerraformModuleSourceVersionConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleSourceVersion(),
					Message: "module 'my_module' source 'git::https://gitlab.example.com/test.git?ref=feature/xxx' [ref='feature/xxx'] does not match any allowed_versions pattern",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 70},
					},
				},
			},
		},
	}

	rule := NewTerraformModuleSourceVersion()
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

const testTerraformModuleSourceVersionConfig = `
rule "terraform_module_source_version" {
  enabled = true
}
`
