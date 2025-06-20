package rules

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
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
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		fileContents, _, diags := file.Body.PartialContent(&hcl.BodySchema{
			Blocks: []hcl.BlockHeaderSchema{
				{
					Type:       "module",
					LabelNames: []string{"name"},
				},
				{
					Type:       "resource",
					LabelNames: []string{"type", "name"},
				},
				{
					Type:       "data",
					LabelNames: []string{"type", "name"},
				},
			},
		})
		if diags.HasErrors() {
			return diags
		}

		for _, fileContent := range fileContents.Blocks {
			var content *hcl.BodyContent
			switch fileContent.Type {
			case "module":
				content, _, diags = fileContent.Body.PartialContent(&hcl.BodySchema{
					Attributes: []hcl.AttributeSchema{
						{Name: "source"},
						{Name: "count"},
						{Name: "for_each"},
						{Name: "providers"},
					},
				})
			case "resource", "data":
				content, _, diags = fileContent.Body.PartialContent(&hcl.BodySchema{
					Attributes: []hcl.AttributeSchema{
						{Name: "count"},
						{Name: "for_each"},
						{Name: "provider"},
					},
					Blocks: []hcl.BlockHeaderSchema{
						{
							Type:       "lifecycle",
							LabelNames: []string{},
						},
					},
				})
			}
			if diags.HasErrors() {
				return diags
			}

			// currentRange is used to check the arrangement of 'source', 'count'/'for_each' and 'providers'.
			// lastAttr is used to check new line after last argument.
			currentRange := fileContent.DefRange
			var lastAttr hcl.Attribute

			// Check meta argument 'source'.
			source, sourceExist := content.Attributes["source"]
			if sourceExist {
				// 'source' is expected to be placed under 'module "my_module" {'.
				if currentRange.End.Line+1 != source.Range.Start.Line {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%s '%s' has invalid 'source' meta argument arrangement", fileContent.Type, strings.Join(fileContent.Labels, ".")),
						source.Range,
					); err != nil {
						return err
					}
					continue
				}
				currentRange = source.Range
				lastAttr = *source
			}

			// Check meta arguments 'count' or 'for_each'.
			// Only one of 'count' and 'for_each' are allowed in Terraform meta argument.
			count, countExist := content.Attributes["count"]
			forEach, forEachExist := content.Attributes["for_each"]
			if countExist {
				checkLine := currentRange.End.Line
				if sourceExist {
					// 'count' is expected to be placed under 'source = "./my-module/"' with an extra newline.
					checkLine += 2
				} else {
					// 'count' is expected to be placed under 'resource "my_type" "my_name" {'.
					checkLine += 1
				}
				if checkLine != count.Range.Start.Line {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%s '%s' has invalid 'count' meta argument arrangement", fileContent.Type, strings.Join(fileContent.Labels, ".")),
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
			} else if forEachExist {
				checkLine := currentRange.End.Line
				if sourceExist {
					// 'for_each' is expected to be placed under 'source = "./my-module/"' with an extra newline.
					checkLine += 2
				} else {
					// 'for_each' is expected to be placed under 'resource "my_type" "my_name" {'.
					checkLine += 1
				}
				if checkLine != forEach.Range.Start.Line {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%s '%s' has invalid 'for_each' meta argument arrangement", fileContent.Type, fileContent.Labels[0]),
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

			// Check meta arguments 'providers' or 'provider'.
			var providerExist bool
			var provider *hcl.Attribute
			var errMsg string

			switch fileContent.Type {
			case "module":
				provider, providerExist = content.Attributes["providers"]
				errMsg = fmt.Sprintf("%s '%s' has invalid 'providers' meta argument arrangement",
					fileContent.Type, strings.Join(fileContent.Labels, "."))
			case "resource", "data": // Combined case since they have identical handling
				provider, providerExist = content.Attributes["provider"]
				errMsg = fmt.Sprintf("%s '%s' has invalid 'provider' meta argument arrangement",
					fileContent.Type, strings.Join(fileContent.Labels, "."))
			}

			if providerExist {
				checkLine := currentRange.End.Line
				if sourceExist || countExist || forEachExist {
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

			// Check new line after 'source', 'count', 'for_each', 'providers' or 'provider.
			if sourceExist || countExist || forEachExist || providerExist {
				lines := bytes.Split(file.Bytes, []byte("\n"))
				checkLine := lines[lastAttr.Range.End.Line]
				if strings.TrimSpace(string(checkLine)) != "" {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%s '%s' has missing new line after meta argument '%s'", fileContent.Type, strings.Join(fileContent.Labels, "."), lastAttr.Name),
						lastAttr.Range,
					); err != nil {
						return err
					}
				}
			}

			// Check meta argument 'lifecycle'.
			lifeCycleBlocks := content.Blocks.OfType("lifecycle")
			if len(lifeCycleBlocks) > 0 {
				var lifecycleFullRange hcl.Range
				var contentEndLine int
				contentFile, err := runner.GetFile(fileContent.DefRange.Filename)
				if err != nil {
					return err
				}

				// Locate the resource to get completed resource range.
				resourceFileBody := contentFile.Body.(*hclsyntax.Body)
				for _, fileBlock := range resourceFileBody.Blocks {
					if fileBlock.Type == "resource" || fileBlock.Type == "data" {
						if fileBlock.Labels[0]+fileBlock.Labels[1] == fileContent.Labels[0]+fileContent.Labels[1] {
							contentEndLine = fileBlock.Range().End.Line
							for _, bodyBlock := range fileBlock.Body.Blocks {
								if bodyBlock.Type == "lifecycle" {
									lifecycleFullRange = bodyBlock.Range()
									break
								}
							}
							break
						}
					}
				}

				if lifecycleFullRange.End.Line+1 != contentEndLine {
					if err := runner.EmitIssue(
						r,
						fmt.Sprintf("%s '%s' has invalid 'lifecycle' meta argument arrangement", fileContent.Type, strings.Join(fileContent.Labels, ".")),
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
