package rules

import (
	"fmt"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-getter"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformModuleSourceVersion struct {
	tflint.DefaultRule
}

// NewTerraformModuleSourceVersion returns a new rule
func NewTerraformModuleSourceVersion() *TerraformModuleSourceVersion {
	return &TerraformModuleSourceVersion{}
}

type TerraformModuleSourceVersionConfig struct {
	AllowedVersions []string `hclext:"allowed_versions,optional"`
}

// Name returns the rule name
func (r *TerraformModuleSourceVersion) Name() string {
	return "terraform_module_source_version"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModuleSourceVersion) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformModuleSourceVersion) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModuleSourceVersion) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_source_version.md"
}

// Link returns the rule reference link
func (r *TerraformModuleSourceVersion) Check(runner tflint.Runner) error {
	config := &TerraformModuleSourceVersionConfig{}

	if err := runner.DecodeRuleConfig(r.Name(), config); err != nil {
		return err
	}

	modules, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "module",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{
							Name: "source",
						},
					},
				},
			},
		},
	}, &tflint.GetModuleContentOption{ExpandMode: tflint.ExpandModeNone})
	if err != nil {
		return err
	}

	for _, module := range modules.Blocks {
		sourceAttr, sourceExist := module.Body.Attributes["source"]
		if !sourceExist {
			continue
		}

		var sourceValue string
		if err := runner.EvaluateExpr(sourceAttr.Expr, &sourceValue, nil); err != nil {
			return err
		}

		source, err := getter.Detect(sourceValue, filepath.Dir(module.DefRange.Filename), []getter.Detector{
			new(getter.GitHubDetector),
			new(getter.GitDetector),
			new(getter.FileDetector),
		})
		if err != nil {
			return err
		}

		u, err := url.ParseRequestURI(source)
		if err != nil {
			if _err := runner.EmitIssue(
				r,
				fmt.Sprintf("module '%s' source '%s' is not a valid URL", module.Labels[0], sourceValue),
				sourceAttr.Expr.Range(),
			); _err != nil {
				return _err
			}
			continue
		}

		// Only enforce version checks for git-based sources
		if u.Scheme != "git" {
			continue
		}

		if u.Opaque != "" {
			query := u.RawQuery
			u, err = url.Parse(strings.TrimPrefix(u.Opaque, ":"))
			if err != nil {
				return err
			}
			u.RawQuery = query
		}

		query := u.Query()

		revision := query.Get("ref")
		key := "ref"

		if revision == "" {
			revision = query.Get("rev")
			key = "rev"

			if revision == "" {
				if _err := runner.EmitIssue(
					r,
					fmt.Sprintf(`module '%s' source '%s' is not pinned (missing ?ref= or ?rev= in the URL).`, module.Labels[0], sourceValue),
					sourceAttr.Expr.Range(),
				); _err != nil {
					return _err
				}
				continue
			}
		}

		_, err = semver.NewVersion(revision)
		if err != nil {
			allowed := false
			for _, version := range config.AllowedVersions {
				re := regexp.MustCompile(version)
				if re.MatchString(revision) {
					allowed = true
					break
				}
			}

			if !allowed {
				if _err := runner.EmitIssue(
					r,
					fmt.Sprintf("module '%s' source '%s' [%s='%s'] does not match any allowed_versions pattern", module.Labels[0], sourceValue, key, revision),
					sourceAttr.Expr.Range(),
				); _err != nil {
					return _err
				}
				continue
			}
		}
	}

	return nil
}
