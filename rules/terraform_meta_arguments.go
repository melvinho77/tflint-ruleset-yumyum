package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformMetaArguments checks the sequences of meta arguments
type TerraformMetaArguments struct {
	tflint.DefaultRule
}

// NewTerraformMetaArguments returns a new rule
func NewTerraformMetaArguments() *TerraformMetaArguments {
	return &TerraformMetaArguments{}
}

// Name returns the rule name
func (r *TerraformMetaArguments) Name() string {
	return "terraform_meta_arguments"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformMetaArguments) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformMetaArguments) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformMetaArguments) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_meta_arguments.md"
}

// Check checks whether variables have type
func (r *TerraformMetaArguments) Check(runner tflint.Runner) error {
	blocks, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "source"},
						{Name: "count"},
						{Name: "for_each"},
						{Name: "providers"},
					},
				},
			},
			{
				Type:       "resource",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "count"},
						{Name: "for_each"},
						{Name: "provider"},
					},
					Blocks: []hclext.BlockSchema{
						{
							Type:       "lifecycle",
							LabelNames: []string{},
							Body:       &hclext.BodySchema{},
						},
					},
				},
			},
			{
				Type:       "data",
				LabelNames: []string{"type", "name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "count"},
						{Name: "for_each"},
						{Name: "provider"},
					},
					Blocks: []hclext.BlockSchema{
						{
							Type:       "lifecycle",
							LabelNames: []string{},
							Body:       &hclext.BodySchema{},
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, block := range blocks.Blocks {
		// currentRange is used to check the arrangement of 'source', 'count'/'for_each' and 'providers'.
		// lastAttr is used to check new line after last argument.
		currentRange := block.DefRange
		var lastAttr hclext.Attribute

		var sourceExists bool
		// 'source' is required in every modules, if the resource type is not module, skip this checking.
		if block.Type == "module" {
			var source *hclext.Attribute
			source, sourceExists = block.Body.Attributes["source"]
			if sourceExists {
				// 'source' is expected to be placed under 'module "my_module" {'.
				if currentRange.End.Line+1 != source.Range.Start.Line {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Invalid 'source' argument arrangement in %s '%s'", block.Type, strings.Join(block.Labels, ".")),
						source.Range,
					); err != nil {
						return err
					}
					continue
				}
				currentRange = source.Range
				lastAttr = *source
			}
		} else {
			sourceExists = false
		}

		// Only one of 'count' and 'for_each' are allowed in Terraform meta argument.
		count, countExists := block.Body.Attributes["count"]
		forEach, forEachExists := block.Body.Attributes["for_each"]
		if countExists {
			checkLine := currentRange.End.Line
			if sourceExists {
				// 'count' is expected to be placed under 'source = "./my-module/"' with an extra newline.
				checkLine += 2
			} else {
				// 'count' is expected to be placed under 'resource "my_type" "my_name" {'.
				checkLine += 1
			}
			if checkLine != count.Range.Start.Line {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Invalid 'count' argument arrangement in %s '%s'", block.Type, strings.Join(block.Labels, ".")),
					count.Range,
				); err != nil {
					return err
				}
				continue
			}
			currentRange = count.Range
			if count.Range.Start.Line > lastAttr.Range.Start.Line {
				lastAttr = *count
			}
		} else if forEachExists {
			checkLine := currentRange.End.Line
			if sourceExists {
				// 'for_each' is expected to be placed under 'source = "./my-module/"' with an extra newline.
				checkLine += 2
			} else {
				// 'for_each' is expected to be placed under 'resource "my_type" "my_name" {'.
				checkLine += 1
			}
			if checkLine != forEach.Range.Start.Line {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Invalid 'for_each' argument arrangement in %s '%s'", block.Type, block.Labels[0]),
					forEach.Range,
				); err != nil {
					return err
				}
				continue
			}
			currentRange = forEach.Range
			if forEach.Range.Start.Line > lastAttr.Range.Start.Line {
				lastAttr = *forEach
			}
		}

		var providerExists bool
		var provider *hclext.Attribute
		var errMsg string
		if block.Type == "module" {
			provider, providerExists = block.Body.Attributes["providers"]
			errMsg = fmt.Sprintf("Invalid 'providers' argument arrangement in %s '%s'", block.Type, strings.Join(block.Labels, "."))
		} else if block.Type == "resource" || block.Type == "data" {
			provider, providerExists = block.Body.Attributes["provider"]
			errMsg = fmt.Sprintf("Invalid 'provider' argument arrangement in %s '%s'", block.Type, strings.Join(block.Labels, "."))
		}
		if providerExists {
			checkLine := currentRange.End.Line
			if sourceExists || countExists || forEachExists {
				// 'providers' is expected to be placed under 'count=0' or 'for_each={}' with an extra newline.
				checkLine += 2
			} else {
				// 'providers' is expected to be placed under 'resource "my_type" "my_name" {'.
				checkLine += 1
			}
			if checkLine != provider.Range.Start.Line {
				if err := runner.EmitIssue(
					r,
					errMsg,
					provider.Range,
				); err != nil {
					return err
				}
				continue
			}
			currentRange = provider.Range
			if provider.Range.Start.Line > lastAttr.Range.Start.Line {
				lastAttr = *provider
			}
		}

		// Check new line after last checked argument if any.
		if sourceExists || countExists || forEachExists || providerExists {
			file, err := runner.GetFile(lastAttr.Range.Filename)
			if err != nil {
				return err
			}
			lines := bytes.Split(file.Bytes, []byte("\n"))
			checkLine := lines[lastAttr.Range.End.Line]
			if strings.TrimSpace(string(checkLine)) != "" {
				if err := runner.EmitIssue(
					r,
					fmt.Sprintf("Missing new line after '%s' in %s '%s'", lastAttr.Name, block.Type, strings.Join(block.Labels, ".")),
					lastAttr.Range,
				); err != nil {
					return err
				}
			}
		}

		// lifecycle block is expected in resource and data source only.
		if block.Type == "resource" || block.Type == "data" {
			var lifecycleFullRange hcl.Range
			lifeCycleBlocks := block.Body.Blocks.OfType("lifecycle")
			if len(lifeCycleBlocks) > 0 {
				var contentEndLine int
				var lifeCycleEndLine int
				contentFile, err := runner.GetFile(block.DefRange.Filename)
				if err != nil {
					return err
				}

				// Locate the resource to get completed resource range.
				resourceFileBody := contentFile.Body.(*hclsyntax.Body)
				for _, fileBlock := range resourceFileBody.Blocks {
					if fileBlock.Type == "resource" || fileBlock.Type == "data" {
						if fileBlock.Labels[0]+fileBlock.Labels[1] == block.Labels[0]+block.Labels[1] {
							contentEndLine = fileBlock.Range().End.Line
							for _, bodyBlock := range fileBlock.Body.Blocks {
								if bodyBlock.Type == "lifecycle" {
									lifeCycleEndLine = bodyBlock.Range().End.Line
									lifecycleFullRange = bodyBlock.Range()
									break
								}
							}
							break
						}
					}
				}

				if lifeCycleEndLine+1 != contentEndLine {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("Invalid 'lifecycle' argument arrangement in %s '%s'", block.Type, strings.Join(block.Labels, ".")),
						lifecycleFullRange,
					); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
