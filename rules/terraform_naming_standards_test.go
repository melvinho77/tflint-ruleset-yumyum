package rules

import (
	"fmt"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformNamingStandards_Variable_DefaultEmpty(t *testing.T) {
	testVariableSnakeCase(t, "default config", "format: snake_case", `
rule "terraform_naming_standards" {
  enabled = true
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultFormat(t *testing.T) {
	testVariableMixedSnakeCase(t, `default config (format="mixed_snake_case")`, `
rule "terraform_naming_standards" {
  enabled = true
  format  = "mixed_snake_case"
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultCustom(t *testing.T) {
	testVariableSnakeCase(t, `default config (custom="^[a-z_]+$")`, "RegExp: ^[a-z][a-z]*(_[a-z]+)*$", `
rule "terraform_naming_standards" {
  enabled = true
  custom  = "^[a-z][a-z]*(_[a-z]+)*$"
}`)
}

// func Test_TerraformNamingStandards_Variable_DefaultDisabled(t *testing.T) {
// 	testVariableDisabled(t, `default config (format=null)`, `
// rule "terraform_naming_standards" {
//   enabled = true
//   format  = "none"
// }`)
// }

func Test_TerraformNamingStandards_Variable_DefaultFormat_OverrideFormat(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "format: snake_case", `
rule "terraform_naming_standards" {
  enabled = true
  format  = "mixed_snake_case"

  variable {
    format = "snake_case"
  }
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultFormat_OverrideCustom(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "RegExp: ^[a-z][a-z]*(_[a-z]+)*$", `
rule "terraform_naming_standards" {
  enabled = true
  format  = "mixed_snake_case"

  variable {
    custom = "^[a-z][a-z]*(_[a-z]+)*$"
  }
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultCustom_OverrideFormat(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "format: snake_case", `
rule "terraform_naming_standards" {
  enabled = true
  custom  = "^ignored$"

  variable {
    format = "snake_case"
  }
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultCustom_OverrideCustom(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "RegExp: ^[a-z][a-z]*(_[a-z]+)*$", `
rule "terraform_naming_standards" {
  enabled = true
  custom  = "^ignored$"

  variable {
    custom = "^[a-z][a-z]*(_[a-z]+)*$"
  }
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultDisabled_OverrideFormat(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "format: snake_case", `
rule "terraform_naming_standards" {
  enabled = true
  format  = "none"

  variable {
    format = "snake_case"
  }
}`)
}

func Test_TerraformNamingStandards_Variable_DefaultDisabled_OverrideCustom(t *testing.T) {
	testVariableSnakeCase(t, `overridden config (format="snake_case")`, "RegExp: ^[a-z][a-z]*(_[a-z]+)*$", `
rule "terraform_naming_standards" {
  enabled = true
  format  = "none"

  variable {
    custom = "^[a-z][a-z]*(_[a-z]+)*$"
  }
}`)
}

// func Test_TerraformNamingStandards_Variable_DefaultEmpty_OverrideDisabled(t *testing.T) {
// 	testVariableDisabled(t, `overridden config (format=null)`, `
// rule "terraform_naming_standards" {
//   enabled = true

//   variable {
//     format = "none"
//   }
// }`)
// }

// func Test_TerraformNamingStandards_Variable_DefaultFormat_OverrideDisabled(t *testing.T) {
// 	testVariableDisabled(t, `overridden config (format=null)`, `
// rule "terraform_naming_standards" {
//   enabled = true
//   format  = "snake_case"

//   variable {
//     format = "none"
//   }
// }`)
// }

func testVariableSnakeCase(t *testing.T, testType, formatName, config string) {
	rule := NewTerraformNamingStandards()

	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: fmt.Sprintf("variable: %s - Invalid snake_case with dash", testType),
			Content: `
variable "dash-name" {
  description = "invalid"
}
		`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `dash-name` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid snake_case with camelCase", testType),
			Content: `
variable "camelCased" {
  description = "invalid"
}
		`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `camelCased` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 22},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid snake_case with double underscore", testType),
			Content: `
variable "foo__bar" {
  description = "invalid"
}
		`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo__bar` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 20},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid snake_case with underscore tail", testType),
			Content: `
variable "foo_bar_" {
  description = "invalid"
}
		`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar_` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 20},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid snake_case with Mixed_Snake_Case", testType),
			Content: `
variable "Foo_Bar" {
  description = "invalid"
}`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `Foo_Bar` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
}`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid single word", testType),
			Content: `
variable "foo" {
  description = "valid"
}`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid object type snake_case.", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = object({
    foo_bar_one = string
    foo_bar_two = string
  })
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid object type with valid sub-objects snake_case.", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = object({
    foo_bar        = string
    foo_bar_config = object({
      foo_bar_option = string
      foo_bar_data   = string
    })
  })
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid object type snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "invalid"
  type = object({
    foo_bar_option = string
    foo_bar_       = string
    foo__bar       = string
  })
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.foo_bar_` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 8, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.foo__bar` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 8, Column: 5},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid object type snake_case with sub-object attributes", testType),
			Content: `
variable "foo_bar" {
  description = "invalid"
  type = object({
    foo_bar_text = string
    foo_bar_bool = string
    fooBarObject = object({
      foo_bar_subText = string
    })
  })
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.fooBarObject` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.fooBarObject.foo_bar_subText` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 10, Column: 5},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid map object type snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = map(object({
    foo_bar_string  = string
    foo_bar_number  = number
    foo_bar_obj     = object({
      foo_bar_nested_one = bool
      foo_bar_nested_two = number
    })
  }))
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid map object type snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "invalid"
  type = map(object({
    foo_bar_string  = string
    foo_barNumber   = number
  }))
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.foo_barNumber` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 6},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid map object type with sub-objects snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "invalid"
  type = map(object({
    foo_bar_string  = string
    foo_barNumber   = number
    fooBar_object   = object({
      foo_bar_nested_one = bool
      fooBar_nestedTwo   = number
    })
  }))
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.foo_barNumber` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.fooBar_object` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: fmt.Sprintf("variable name `foo_bar.fooBar_object.fooBar_nestedTwo` must match the following %s", formatName),
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
			},
		},
	}

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

func testVariableMixedSnakeCase(t *testing.T, testType, config string) {
	rule := NewTerraformNamingStandards()

	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: fmt.Sprintf("variable: %s - Invalid mixed_snake_case with dash", testType),
			Content: `
variable "dash-name" {
  description = "invalid"
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable name `dash-name` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 21},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid mixed_snake_case with double underscore", testType),
			Content: `
variable "Foo__Bar" {
  description = "invalid"
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable name `Foo__Bar` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 20},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid mixed_snake_case with underscore tail", testType),
			Content: `
variable "Foo_Bar_" {
  description = "invalid"
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable name `Foo_Bar_` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 20},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid single word", testType),
			Content: `
variable "foo" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid Mixed_Snake_Case", testType),
			Content: `
variable "Foo_Bar" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid single word with upper characters", testType),
			Content: `
variable "foo" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid PascalCase", testType),
			Content: `
variable "PascalCase" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid camelCase", testType),
			Content: `
variable "camelCase" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valud object type with mixed_snake_case, camelCase and PascalCase", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = object({
    snake_case = string
    camelCase  = number
    PascalCase = bool
  })
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid object type with mixed_snake_case with dash", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = object({
    foo-bar-one = string
    foo-bar-two = number
  })
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable name `foo_bar.foo-bar-one` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable name `foo_bar.foo-bar-two` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 7, Column: 5},
					},
				},
			},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid object type with sub-objects with mixed_snake_case, PascalCase, camelCase", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = object({
    foo_bar_one     = string
    foo_bar_two     = number
    foo_bar_subObj = object({
      mixed_snake_case = string
      camelCase        = string
      PascalCase       = object({
        foo_bar_test = string
      })
    })
  })
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid map object type with mixed_snake_case, PascalCase, camelCase", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = map(object({
    mixed_snake_case = string
    PascalCase       = number
    camelCase        = bool
  }))
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid map object type with sub-objects with mixed_snake_case, PascalCase, camelCase", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = map(object({
    mixed_snake_case = string
    PascalCase       = number
    camelCase        = object({
      mixed_snake_case = string
      PascalCase       = number
      camelCase        = bool
    })
  }))
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Invalid map object type with sub-objects with mixed_snake_case with dash", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
  type = map(object({
    mixed-snake-case-dash = string
    test_object           = object({
      mixed-snake-case-dash = string
      PascalCase            = number
      camelCase             = bool
    })
  }))
}
`,
			Config: config,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable name `foo_bar.mixed-snake-case-dash` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable name `foo_bar.test_object.mixed-snake-case-dash` must match the following format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 4, Column: 3},
						End:      hcl.Pos{Line: 11, Column: 6},
					},
				},
			},
		},
	}

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

func testVariableDisabled(t *testing.T, testType, config string) {
	rule := NewTerraformNamingStandards()

	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		{
			Name: fmt.Sprintf("variable: %s - Valid mixed_snake_case with dash", testType),
			Content: `
variable "dash-name" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid snake_case", testType),
			Content: `
variable "foo_bar" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid single word", testType),
			Content: `
variable "foo" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid Mixed_Snake_Case", testType),
			Content: `
variable "Foo_Bar" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid single word upper characters", testType),
			Content: `
variable "Foo" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid PascalCase", testType),
			Content: `
variable "PascalCase" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
		{
			Name: fmt.Sprintf("variable: %s - Valid camelCase", testType),
			Content: `
variable "camelCase" {
  description = "valid"
}
`,
			Config:   config,
			Expected: helper.Issues{},
		},
	}

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
