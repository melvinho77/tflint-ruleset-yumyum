package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformMetaArguments(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "source only in module",
			Content: `
module "my_module" {
  source = "./my-module/"

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "source and count only in module",
			Content: `
module "my_module" {
  source = "./my-module/"

  count = 3

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "source and for_each only in module",
			Content: `
module "my_module" {
  source = "./my-module/"

  for_each = {}

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "source, count and providers in module",
			Content: `
module "my_module" {
  source = "./my-module/"

  count = {}

  providers = {}

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "source, for_each and providers in module",
			Content: `
module "my_module" {
  source = "./my-module/"

  for_each = {}

  providers = {}

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count only in resource",
			Content: `
resource "foo" "my_resource" {
  count = 3

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "for_each only in resource",
			Content: `
resource "foo" "my_resource" {
  for_each = {}

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "provider only in resource",
			Content: `
resource "foo" "my_resource" {
  provider = foo.default

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count and provider in resource",
			Content: `
resource "foo" "my_resource" {
  count = 3

  provider = foo.default

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count and provider in data source",
			Content: `
data "foo" "my_resource" {
  count = 3

  provider = foo.default

  name = "my name"
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "lifecycle in resource",
			Content: `
resource "foo" "my_resource" {
  name = "my name"

  lifecycle {}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "count, provider and lifecycle in resource",
			Content: `
resource "foo" "my_resource" {
  count = 3

  provider = foo.default

  name = "my name"

  lifecycle {}
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "source and count in module, invalid arrangement",
			Content: `
module "my_module" {
  count = 3

  source = "./my-module/"

  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "module 'my_module' has invalid 'source' meta argument arrangement",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 26},
					},
				},
			},
		},
		{
			Name: "count, provider in resource, invalid arrangement",
			Content: `
resource "foo" "my_resource" {
  provider = foo.default

  count = 3

  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "resource 'foo.my_resource' has invalid 'count' meta argument arrangement",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 12},
					},
				},
			},
		},
		{
			Name: "count, provider and lifecycle in resource, invalid arrangement",
			Content: `
resource "foo" "my_resource" {
  count = 3

  provider = foo.default

  lifecycle {}

  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "resource 'foo.my_resource' has invalid 'lifecycle' meta argument arrangement",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 15},
					},
				},
			},
		},
		{
			Name: "source, count and providers in module, invalid arrangement",
			Content: `
module "my_module" {
  source = "./my-module/"

  providers = {}

  count = 3

  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "module 'my_module' has invalid 'count' meta argument arrangement",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 12},
					},
				},
			},
		},
		{
			Name: "source only in module, missing new line",
			Content: `
module "my_module" {
  source = "./my-module/"
  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "module 'my_module' has missing new line after meta argument 'source'",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 3, Column: 26},
					},
				},
			},
		},
		{
			Name: "source and count only in module, missing new line",
			Content: `
module "my_module" {
  source = "./my-module/"

  count = 3
  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "module 'my_module' has missing new line after meta argument 'count'",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 3},
						End:      hcl.Pos{Line: 5, Column: 12},
					},
				},
			},
		},
		{
			Name: "source, count and providers in module, missing new line",
			Content: `
module "my_module" {
  source = "./my-module/"

  count = 3

  providers = {}
  name = "my name"
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformMetaArguments(),
					Message: "module 'my_module' has missing new line after meta argument 'providers'",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 17},
					},
				},
			},
		},
	}

	rule := NewTerraformMetaArguments()
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"main.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
