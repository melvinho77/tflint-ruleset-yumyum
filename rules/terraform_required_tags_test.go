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
			Name: "resource with the correct tags, but terraform module with no local variable `tags`.",
			Content: `
resource "my_resource" "my_resource_name" {
  tags = {
    my_required_tag = "my_tag"
  }
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "missing required local variable `tags`",
					Range: hcl.Range{
						Start: hcl.Pos{Line: 1, Column: 1},
						End:   hcl.Pos{Line: 1, Column: 1},
					},
				},
			},
		},
		{
			Name: "resource with the incorrect tags, but terraform module with no local variable `tags`.",
			Content: `
resource "my_resource" "my_resource_name" {
  tags = {
    my_incorrect_tag = "my_tag"
  }
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "missing required local variable `tags`",
					Range: hcl.Range{
						Start: hcl.Pos{Line: 1, Column: 1},
						End:   hcl.Pos{Line: 1, Column: 1},
					},
				},
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "my_resource 'my_resource_name' is missing required tags: [my_required_tag]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 10},
						End:      hcl.Pos{Line: 5, Column: 4},
					},
				},
			},
		},
		{
			Name: "resource without any tags block, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

resource "my_resource" "my_resource_name" {
  name = "test"
}
`,
			Config:   testTerraformRequiredTagsConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "resource with the correct tags, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

resource "my_resource" "my_resource_name" {
  name = "test"

  tags = {
    my_required_tag = "my_tag"
  }
}
`,
			Config: `
rule "terraform_required_tags" {
  enabled            = true

  tags               = ["my_required_tag"]
  excluded_resources = ["my_excluded_resource"]
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource with the incorrect tags, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

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
						Start:    hcl.Pos{Line: 11, Column: 10},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
			},
		},
		{
			Name: "resource that is excluded from the rule, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

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
			Name: "aws resource with the correct tags, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

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
			Name: "aws resource with but required tags but no `Name` tag, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

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
						Start:    hcl.Pos{Line: 11, Column: 10},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
			},
		},
		{
			Name: "aws resource with the incorrect tags and no `Name` tags, and local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    my_required_tag = "my_tag"
  }
}

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
					Message: "aws_resource 'my_resource_name' is missing required tags: [my_required_tag]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 11, Column: 10},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'my_resource_name' is missing required tag: Name",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 11, Column: 10},
						End:      hcl.Pos{Line: 13, Column: 4},
					},
				},
			},
		},
		{
			Name: "aws resource with `Name` tag, merged with local variable `tags` with the correct keys.",
			Content: `
locals {
  tags = {
    brand           = "foo"
    project         = "bar"
    env             = "dev"
    my_required_tag = "test"
  }
}

resource "aws_instance" "aws_resource_name" {
  tags = merge(local.tags, {
    Name = "aws_instance-test"
  })
}
`,
			Config:   testTerraformRequiredTagsConfig,
			Expected: helper.Issues{},
		},
		{
			Name: "aws resource with no `Name` tag, merged with local variable `tags` with the incorrect keys.",
			Content: `
locals {
  tags = {
    foo = "bar"
    bar = "foo"
  }
}

resource "aws_resource" "aws_resource_name" {
  tags = merge(local.tags, {
    testTag = "test"
  })
}
`,
			Config: testTerraformRequiredTagsConfig,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'aws_resource_name' is missing required tag: Name",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 10, Column: 10},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "aws_resource 'aws_resource_name' is missing required tags: [my_required_tag]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 10, Column: 10},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
			},
		},
		{
			Name: "resource using local variable as tags, and local variable `tags` with the correct default keys.",
			Content: `
locals {
  tags = {
    brand                = "foo"
    env                  = "bar"
    project              = "dev"
    devops_project_kind  = "foo"
    devops_project_group = "bar"
    devops_project_name  = "dev"
  }
}

resource "my_resource" "my_resource_name" {
  tags = local.tags
}
`,
			Config: `
rule "terraform_required_tags" {
  enabled = true
}
`,
			Expected: helper.Issues{},
		},
		{
			Name: "resource using local variable as tags, and local variable `tags` with the incorrect default keys.",
			Content: `
locals {
  tags = {
    brand                = "foo"
    env                  = "bar"
    project              = "dev"
  }
}

resource "my_resource" "my_resource_name" {
  tags = local.tags
}
`,
			Config: `
rule "terraform_required_tags" {
  enabled = true
}
`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformRequiredTags(),
					Message: "my_resource 'my_resource_name' is missing required tags: [devops_project_kind, devops_project_group, devops_project_name]",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 11, Column: 10},
						End:      hcl.Pos{Line: 11, Column: 20},
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
