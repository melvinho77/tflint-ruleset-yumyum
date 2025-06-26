package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	// Set default required tags if none are specified
	if len(config.Tags) == 0 {
		config.Tags = []string{
			"brand",
			"env",
			"project",
			"devops_project_kind",
			"devops_project_group",
			"devops_project_name",
		}
	}

	// Get and evaluate `local.tags`
	localTagsAttr, err := r.getLocalTags(runner)
	if err != nil {
		return err
	}

	var localTags cty.Value
	var localTagKeys []string

	if localTagsAttr == nil {
		err := runner.EmitIssue(
			r,
			"missing required local variable `tags`",
			hcl.Range{
				Start: hcl.Pos{Line: 1, Column: 1},
				End:   hcl.Pos{Line: 1, Column: 1},
			},
		)
		if err != nil {
			return err
		}
	} else {
		err := runner.EvaluateExpr(localTagsAttr.Expr, func(val cty.Value) error {
			localTags = val
			return nil
		}, nil)
		if err != nil {
			return err
		}

		// Store all the key from `local.tags` if applicable
		if !localTags.IsKnown() || localTags.IsNull() || !localTags.CanIterateElements() {
			return nil
		} else {
			for it := localTags.ElementIterator(); it.Next(); {
				k, _ := it.Element()
				localTagKeys = append(localTagKeys, k.AsString())
			}
		}
	}

	// Parse resources and check their `tags` blocks
	resources, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "tags"},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		if slices.Contains(config.ExcludedResources, resource.Labels[0]) {
			continue
		}

		tagsAttr, tagsExist := resource.Body.Attributes["tags"]
		if !tagsExist {
			continue
		}

		var tagKeys []string
		var evalErr error

		switch expr := tagsAttr.Expr.(type) {
		// Usage of function calls like merge(local.tags, { ... })
		case *hclsyntax.FunctionCallExpr:
			if expr.Name == "merge" {
				for _, arg := range expr.Args {
					// If the argument is `local.tags`, inject keys directly
					if varExpr, ok := arg.(*hclsyntax.ScopeTraversalExpr); ok && varExpr.Traversal.RootName() == "local" {
						tagKeys = append(tagKeys, localTagKeys...)
						continue
					}

					// Otherwise, evaluate and extract keys
					err := runner.EvaluateExpr(arg, func(val cty.Value) error {
						if val.IsKnown() && val.CanIterateElements() {
							for it := val.ElementIterator(); it.Next(); {
								k, _ := it.Element()
								tagKeys = append(tagKeys, k.AsString())
							}
						}
						return nil
					}, nil)

					if err != nil {
						evalErr = err
						break
					}
				}
			} else {
				// Ignore any terraform function calls other than `merge`
				continue
			}

		// Direct use of local variable on tags
		// E.g. tags = local.tags
		case *hclsyntax.ScopeTraversalExpr:
			if expr.Traversal.RootName() == "local" {
				tagKeys = append(tagKeys, localTagKeys...)
			} else {
				continue
			}

		default:
			evalErr = runner.EvaluateExpr(tagsAttr.Expr, func(val cty.Value) error {
				if val.IsKnown() && val.CanIterateElements() {
					for it := val.ElementIterator(); it.Next(); {
						k, _ := it.Element()
						tagKeys = append(tagKeys, k.AsString())
					}
				}
				return nil
			}, nil)
		}

		if evalErr != nil {
			return evalErr
		}

		// Remove any duplicated keys if any
		tagKeys = slices.Compact(tagKeys)

		var missing []string

		for _, requiredTags := range config.Tags {
			if !slices.Contains(tagKeys, requiredTags) {
				missing = append(missing, requiredTags)
			}
		}

		// Output linting error if any missing tags are present
		if len(missing) > 0 {
			err := runner.EmitIssue(
				r,
				fmt.Sprintf("%s '%s' is missing required tags: [%s]", resource.Labels[0], resource.Labels[1], strings.Join(missing, ", ")),
				tagsAttr.Expr.Range(),
			)
			if err != nil {
				return err
			}
		}

		// If resource is AWS Cloud resource, check if `Name` tag key exists
		if r.isAwsResource(resource.Labels[0]) && !slices.Contains(tagKeys, "Name") {
			err := runner.EmitIssue(
				r,
				fmt.Sprintf("%s '%s' is missing required tag: Name", resource.Labels[0], resource.Labels[1]),
				tagsAttr.Expr.Range(),
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Function to determine whether resource has `aws_` prefix
func (r *TerraformRequiredTags) isAwsResource(resource string) bool {
	return strings.HasPrefix(resource, "aws_")
}

// Function to determine whether local variable `tags` exists.
func (r *TerraformRequiredTags) getLocalTags(runner tflint.Runner) (*hclext.Attribute, error) {
	locals, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type: "locals",
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "tags"},
					},
				},
			},
		},
	}, nil)

	if err != nil {
		return nil, err
	}

	for _, block := range locals.Blocks {
		if attr, ok := block.Body.Attributes["tags"]; ok {
			return attr, nil
		}
	}

	return nil, nil
}
