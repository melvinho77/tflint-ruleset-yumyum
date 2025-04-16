package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformRequiredTags(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: "resource without any tags block.",
			Content: `
resource "my_resource" "my_resource_name" {
  name = "test"
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource with the correct tags.",
			Content: `
resource "my_resource" "my_resource_name" {
  name = "test"

  tags = {
    my_required_tag = "my_tag"
  }
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource with the incorrect tags.",
			Content: `
resource "my_resource" "my_resource_name" {
  name = "test"

  tags = {
    my_incorrect_tag = "my_tag"
  }
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "my_resource 'my_resource_name' is missing required tags: [my_required_tag]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "resource that is excluded from the rule.",
			Content: `
resource "my_excluded_resource" "my_excluded_resource_name" {
  name = "test"

  tags = {
    my_incorrect_tag = "my_tag"
  }
}
`,
			Config:   testTerraformRequiredTagsConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "aws resource with the correct tags.",
			Content: `
resource "aws_resource" "my_resource_name" {
  name = "test"

  tags = {
    my_required_tag = "my_tag"
    Name            = "test"
  }
}
`,
			Config:   testTerraformRequiredTagsConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "aws resource with but required tags but no `Name` tag.",
			Content: `
resource "aws_resource" "my_resource_name" {
  name = "test"

  tags = {
    my_required_tag = "my_tag"
  }
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'my_resource_name' is missing required tag: Name",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
		{
			Name: "aws resource with the incorrect tags and no `Name` tags.",
			Content: `
resource "aws_resource" "my_resource_name" {
  name = "test"

  tags = {
    not_required_tag = "my_tag"
  }
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'my_resource_name' is missing required tag: Name",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'my_resource_name' is missing required tags: [my_required_tag]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 5, Column: 10},
						End:      hcl.Pos{Line: 7, Column: 4},
					},
				},
			},
		},
	}

	rule := NewTerraformRequiredTags()
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

const testTerraformRequiredTagsConfig = `
rule "terraform_required_tags" {
  enabled            = true

  tags               = ["my_required_tag"]
  excluded_resources = ["my_excluded_resource"]
}
`
