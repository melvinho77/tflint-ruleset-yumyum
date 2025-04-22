package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/zclconf/go-cty/cty"
)

type TerraformRequiredTags struct {
	tflint.DefaultRule
}

// NewTerraformRequiredTags returns a new rule
func NewTerraformRequiredTags() *TerraformRequiredTags {
	return &TerraformRequiredTags{}
}

type terraformRequiredTagsConfig struct {
	Tags              []string `hclext:"tags,optional"`
	ExcludedResources []string `hclext:"excluded_resources,optional"`
}

// Name returns the rule name
func (r *TerraformRequiredTags) Name() string {
	return "terraform_required_tags"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformRequiredTags) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformRequiredTags) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformRequiredTags) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_tags.md"
}

// Check checks whether resources have the required tags if applicable
func (r *TerraformRequiredTags) Check(runner tflint.Runner) error {
	config := &terraformRequiredTagsConfig{}

	fmt.Printf("DEBUG - config.Tags: %#v (len: %d)\n", config.Tags, len(config.Tags))

	if len(config.Tags) == 0 {
		config.Tags = append(config.Tags, []string{"brand", "env", "project", "devops_project_kind", "devops_project_group", "devops_project_name"}...)
	}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	resources, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name: "tags",
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		// SKip this check if the resource type matches any of the string in excluded_resources.
		if slices.Contains(config.ExcludedResources, resource.Labels[0]) {
			continue
		}

		tagsAttr, tagsExist := resource.Body.Attributes["tags"]
		if !tagsExist {
			continue
		}

		err := runner.EvaluateExpr(tagsAttr.Expr, func(val cty.Value) error {
			if !val.IsKnown() || val.IsNull() || !val.CanIterateElements() {
				return nil
			}

			var tagKeys []string
			for it := val.ElementIterator(); it.Next(); {
				k, _ := it.Element()
				tagKeys = append(tagKeys, k.AsString())
			}

			var missing []string
			for _, required := range config.Tags {
				if !slices.Contains(tagKeys, required) {
					missing = append(missing, required)
				}
			}

			if len(missing) > 0 {
				runner.EmitIssue(
					r,
					fmt.Sprintf("%s '%s' is missing required tags: [%s]", resource.Labels[0], resource.Labels[1], strings.Join(missing, ", ")),
					tagsAttr.Expr.Range(),
				)
			}

			if r.isAwsResource(resource.Labels[0]) && !slices.Contains(tagKeys, "Name") {
				runner.EmitIssue(
					r,
					fmt.Sprintf("%s '%s' is missing required tag: Name", resource.Labels[0], resource.Labels[1]),
					tagsAttr.Expr.Range(),
				)
			}

			return nil
		}, nil)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TerraformRequiredTags) isAwsResource(resource string) bool {
	return strings.HasPrefix(resource, "aws_")
}
