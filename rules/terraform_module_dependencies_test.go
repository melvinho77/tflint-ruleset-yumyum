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
}`,
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
					Rule:    NewTerraformModuleDependencies(),
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
			Config: `
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "flexible"
  default_branches = ["main", "master", "default"]
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleDependencies(),
					Message: "module 'my_module' source 'git::https://gitlab.example.com/test/test-module.git?ref=main' uses a default branch as ref (main)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 75},
					},
				},
			},
		},
		{
			Name: "git module is pinned.",
			Content: `
module "my_module" {
  source = "git://gitlab.exapmle.com/test.git?ref=pinned"
  name   = "my_name"
}`,
			Config: `
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "flexible"
  default_branches = ["main", "master", "default"]
}
`,
			Expected: helper.Issues{},
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
					Rule:    NewTerraformModuleDependencies(),
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
			Name: "git module reference is pinned, but style is semver.",
			Content: `
module "my_module" {
  source = "git://gitlab.example.com/test.git?ref=pinned"
  name   = "my_name"
}`,
			Config: `
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "semver"
  default_branches = []
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformModuleDependencies(),
					Message: "module 'my_module' source 'git://gitlab.example.com/test.git?ref=pinned' uses a ref which is not a semantic version string",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 3, Column: 58},
					},
				},
			},
		},
		{
			Name: "git module reference is pinned to semver.",
			Content: `
module "my_module" {
  source = "git://gitlab.example.com/test.git?ref=v1.2.3"
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
	}

	rule := NewTerraformModuleDependencies()
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
