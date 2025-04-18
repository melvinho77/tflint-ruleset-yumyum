package rules

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/hashicorp/go-getter"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type TerraformModuleDependencies struct {
	tflint.DefaultRule
}

// NewTerraformModuleDependecies returns a new rule
func NewTerraformModuleDependencies() *TerraformModuleDependencies {
	return &TerraformModuleDependencies{}
}

type terraformModuleDependenciesConfig struct {
	Style           string   `hclext:"style,optional"`
	DefaultBranches []string `hclext:"default_branches,optional"`
}

// Name returns the rule name
func (r *TerraformModuleDependencies) Name() string {
	return "terraform_module_dependencies"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModuleDependencies) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *TerraformModuleDependencies) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModuleDependencies) Link() string {
	return "https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_dependencies.md"
}

// Link returns the rule reference link
func (r *TerraformModuleDependencies) Check(runner tflint.Runner) error {
	config := &terraformModuleDependenciesConfig{}

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
			new(getter.BitBucketDetector),
			new(getter.GCSDetector),
			new(getter.S3Detector),
			new(getter.FileDetector),
		})
		if err != nil {
			return err
		}

		u, err := url.Parse(source)
		if err != nil {
			return err
		}

		switch u.Scheme {
		case "git", "hg":
		default:
			continue
		}

		if u.Opaque != "" {
			// for git:: or hg:: pseudo-URLs, Opaque is :https, but query will still be parsed
			query := u.RawQuery
			u, err = url.Parse(strings.TrimPrefix(u.Opaque, ":"))
			if err != nil {
				return err
			}

			u.RawQuery = query
		}

		if u.Hostname() == "" {
			_ = runner.EmitIssue(
				r,
				fmt.Sprintf("module '%s' source '%s' is not a valid URL", module.Labels[0], sourceValue),
				sourceAttr.Expr.Range(),
			)
			continue
		}

		query := u.Query()

		// Extract both ref and rev parameters from the source URL
		// Default will be using ref
		ref := query.Get("ref")
		rev := query.Get("rev")

		// Default will be using ref, e.g.
		// Source url: git::https://gitlab.example.com/test/test-module.git?ref=v0.0.1
		// key = "ref", revision = "v0.0.1"
		key := "ref"
		revision := ref

		// If revision query parameter is not empty, it will use rev instead
		if rev != "" {
			key = "rev"
			revision = rev
		}

		if revision == "" {
			_ = runner.EmitIssue(
				r,
				fmt.Sprintf(`module '%s' source '%s' is not pinned (missing ?ref= or ?rev= in the URL).`, module.Labels[0], sourceValue),
				sourceAttr.Expr.Range(),
			)
			continue
		}

		switch config.Style {
		case "flexible":
			// The `flexible` Style only requires a revision that is NOT a default branch, e.g.
			// Valid source   : git::https://gitlab.example.com/test/test-module.git?ref=v1.2
			// Invalid source : git::https://gitlab.example.com/test/test-module.git?ref=main
			for _, branch := range config.DefaultBranches {
				if revision == branch {
					_ = runner.EmitIssue(
						r,
						fmt.Sprintf("module '%s' source '%s' uses a default branch as %s (%s)", module.Labels[0], sourceValue, key, branch),
						sourceAttr.Expr.Range(),
					)
				}
			}
		case "semver":
			// The `semver` Style restricts a revision that is a semantic version, e.g.
			// Valid source   : git::https://gitlab.example.com/test/test-module.git?ref=v0.0.1
			// Invalid source : git::https://gitlab.example.com/test/test-module.git?ref=feature/6969
			_, err := semver.NewVersion(revision)
			if err != nil {
				_ = runner.EmitIssue(
					r,
					fmt.Sprintf("module '%s' source '%s' uses a %s which is not a semantic version string", module.Labels[0], sourceValue, key),
					sourceAttr.Expr.Range(),
				)
			}
		default:
			return fmt.Errorf("invalid Style '%s', expected 'semver' or 'flexible'", config.Style)
		}
	}

	return nil
}
