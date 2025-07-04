package main

import (
	"github.com/styumyum/tflint-ruleset-yumyum/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "yumyum",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewTerraformMetaArguments(),
				rules.NewTerraformAnyTypeVariables(),
				rules.NewTerraformRequiredTags(),
				rules.NewTerraformModuleSourceVersion(),
				rules.NewTerraformVarsObjectKeysNamingConventions(),
				rules.NewTerraformRequiredVariables(),
			},
		},
	})
}
