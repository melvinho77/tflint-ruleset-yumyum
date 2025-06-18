package rules

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformVarsObjectKeysNamingConventions(t *testing.T) {
	rule := NewTerraformVarsObjectKeysNamingConventions()

	tests := []struct {
		Name     string
		Content  string
		Config   string
		Expected helper.Issues
	}{
		// Test cases for `snake_case`
		{
			Name: "valid primitive type variable (snake_case)",
			Content: `
variable "foo" {
  type        = string
  description = "valid."
}

variable "foo_bar" {
  type        = number
  description = "valid."
}

variable "foo_bar_test" {
  type        = boolean
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid primitive type variable (snake_case)",
			Content: `
variable "fooBar" {
  type        = string
  description = "invalid."
}

variable "FooBar" {
  type        = number
  description = "invalid."
}

variable "foo_BarTest" {
  type        = boolean
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `fooBar` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
				{
					Rule:    rule,
					Message: "variable `FooBar` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 18},
					},
				},
				{
					Rule:    rule,
					Message: "variable `foo_BarTest` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 12, Column: 1},
						End:      hcl.Pos{Line: 12, Column: 23},
					},
				},
			},
		},
		{
			Name: "valid complex type - object variable (snake_case)",
			Content: `
variable "foo_bar" {
  type        = object({
    id   = string
    name = string
  })
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - object variable (camelCase name and keys)",
			Content: `
variable "fooBar" {
  type = object({
    idNumber  = string
    userName  = string
  })
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `fooBar` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBar` path `fooBar.idNumber` - attribute `idNumber` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBar` path `fooBar.userName` - attribute `userName` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
			},
		},
		{
			Name: "valid complex type - nested object with sub-objects variable (snake_case)",
			Content: `
variable "foo_bar" {
  type = object({
    settings = object({
      retries    = number
      timeout_ms = number
    })
    metadata = object({
      created_by = string
      tags       = list(string)
    })
  })
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - nested object with sub-objects variable (snake_case)",
			Content: `
variable "fooBar" {
  type = object({
    settings = object({
      retries    = number
      timeoutMs  = number
    })
    metadata = object({
      createdBy = string
      tagList   = list(string)
    })
  })
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `fooBar` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 18},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBar` path `fooBar.metadata.createdBy` - attribute `createdBy` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBar` path `fooBar.metadata.tagList` - attribute `tagList` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBar` path `fooBar.settings.timeoutMs` - attribute `timeoutMs` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 12, Column: 5},
					},
				},
			},
		},
		{
			Name: "valid complex type - map of object with list of object inside (snake_case)",
			Content: `
variable "env_services" {
  type = map(object({
    service_name = string
    endpoints = list(object({
      url         = string
      is_secure   = bool
    }))
  }))
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - map of object with list of object inside (snake_case)",
			Content: `
variable "envServices" {
  type = map(object({
    serviceName = string
    endpoints = list(object({
      endpointURL = string
      isSecure    = bool
    }))
  }))
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `envServices` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 23},
					},
				},
				{
					Rule:    rule,
					Message: "variable `envServices` path `envServices.serviceName` - attribute `serviceName` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `envServices` path `envServices.endpoints.endpointURL` - attribute `endpointURL` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `envServices` path `envServices.endpoints.isSecure` - attribute `isSecure` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 9, Column: 6},
					},
				},
			},
		},
		{
			Name: "valid complex type - deeply nested map of object with list of object with map(object) (snake_case)",
			Content: `
variable "complex_config" {
  type = map(object({
    config_name = string
    rules = list(object({
      rule_type = string
      targets = map(object({
	target_id = string
	metadata  = object({
	  created_at = string
	  owner_id   = string
	})
      }))
    }))
  }))
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - deeply nested map of object with list of object with map(object) (camelCase)",
			Content: `
variable "complexConfig" {
  type = map(object({
    configName = string
    rules = list(object({
      ruleType = string
      targets = map(object({
	targetId = string
	metadata = object({
	  createdAt = string
	  ownerId   = string
	})
      }))
    }))
  }))
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_snakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `complexConfig` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 25},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complexConfig` path `complexConfig.configName` - attribute `configName` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complexConfig` path `complexConfig.rules.ruleType` - attribute `ruleType` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complexConfig` path `complexConfig.rules.targets.targetId` - attribute `targetId` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complexConfig` path `complexConfig.rules.targets.metadata.createdAt` - attribute `createdAt` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complexConfig` path `complexConfig.rules.targets.metadata.ownerId` - attribute `ownerId` must match the following predefined_format: snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
			},
		},
		// Test cases for `mixed_snake_case`
		{
			Name: "valid primitive type variable (mixed_snake_case)",
			Content: `
variable "Foo" {
  type        = string
  description = "valid."
}

variable "Foo1_bar" {
  type        = number
  description = "valid."
}

variable "myVar2_testValue" {
  type        = boolean
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid primitive type variable (mixed_snake_case)",
			Content: `
variable "_foo" {
  type        = string
  description = "invalid."
}

variable "1bar_test" {
  type        = number
  description = "invalid."
}

variable "foo__bar" {
  type        = boolean
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `_foo` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 16},
					},
				},
				{
					Rule:    rule,
					Message: "variable `1bar_test` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 21},
					},
				},
				{
					Rule:    rule,
					Message: "variable `foo__bar` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 12, Column: 1},
						End:      hcl.Pos{Line: 12, Column: 20},
					},
				},
			},
		},
		{
			Name: "valid complex type - object variable (mixed_snake_case)",
			Content: `
variable "Foo_bar" {
  type = object({
    Id1      = string
    User_2   = string
  })
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - object variable (invalid mixed_snake_case name and keys)",
			Content: `
variable "_foo" {
  type = object({
    user__name = string
    id_123_    = string
  })
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `_foo` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 16},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_foo` path `_foo.id_123_` - attribute `id_123_` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_foo` path `_foo.user__name` - attribute `user__name` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
			},
		},
		{
			Name: "valid complex type - deeply nested map of object with list of object with map(object) (mixed_snake_case)",
			Content: `
variable "ComplexConfig1" {
  type = map(object({
    ConfigName123 = string
    Rule_List = list(object({
      RuleType = string
      TargetMap_2 = map(object({
	TargetID_2 = string
	MetaData1  = object({
	  CreatedAt = string
	  OwnerID123   = string
	})
      }))
    }))
  }))
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - deeply nested map of object with list of object with map(object) (invalid mixed_snake_case)",
			Content: `
variable "_ComplexConfig1" {
  type = map(object({
    Config__Name123 = string
    RuleList = list(object({
      RuleType1 = string
      TargetMap__2_ = map(object({
	TargetID__ = string
	MetaData1  = object({
	  created_at = string
	  _OwnerID = string
	})
      }))
    }))
  }))
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `_ComplexConfig1` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 27},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_ComplexConfig1` path `_ComplexConfig1.Config__Name123` - attribute `Config__Name123` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_ComplexConfig1` path `_ComplexConfig1.RuleList.TargetMap__2_.MetaData1._OwnerID` - attribute `_OwnerID` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_ComplexConfig1` path `_ComplexConfig1.RuleList.TargetMap__2_.TargetID__` - attribute `TargetID__` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
				{
					Rule:    rule,
					Message: "variable `_ComplexConfig1` path `_ComplexConfig1.RuleList.TargetMap__2_` - attribute `TargetMap__2_` must match the following predefined_format: mixed_snake_case",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 15, Column: 6},
					},
				},
			},
		},
		// Test cases for `custom_format` - PascalCase
		{
			Name: "valid primitive type variable (PascalCase)",
			Content: `
variable "Foo" {
  type        = string
  description = "valid."
}

variable "FooBar" {
  type        = number
  description = "valid."
}

variable "FooBarTest" {
  type        = boolean
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid primitive type variable (PascalCase)",
			Content: `
variable "foo" {
  type        = string
  description = "invalid."
}

variable "foo_bar" {
  type        = number
  description = "invalid."
}

variable "fooBarTest" {
  type        = boolean
  description = "invalid."
}
		`,
			Config: testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `foo` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 15},
					},
				},
				{
					Rule:    rule,
					Message: "variable `foo_bar` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 7, Column: 1},
						End:      hcl.Pos{Line: 7, Column: 19},
					},
				},
				{
					Rule:    rule,
					Message: "variable `fooBarTest` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 12, Column: 1},
						End:      hcl.Pos{Line: 12, Column: 22},
					},
				},
			},
		},
		{
			Name: "valid complex type - object variable (PascalCase)",
			Content: `
variable "FooBar" {
  type = object({
    Id        = string
    UserName  = string
  })
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - object variable (non-PascalCase name and keys)",
			Content: `
variable "foo_bar" {
  type = object({
    id_number = string
    user_name = string
  })
  description = "invalid."
}
`,
			Config: testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `foo_bar` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 1},
						End:      hcl.Pos{Line: 2, Column: 19},
					},
				},
				{
					Rule:    rule,
					Message: "variable `foo_bar` path `foo_bar.id_number` - attribute `id_number` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
				{
					Rule:    rule,
					Message: "variable `foo_bar` path `foo_bar.user_name` - attribute `user_name` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 3},
						End:      hcl.Pos{Line: 6, Column: 5},
					},
				},
			},
		},
		{
			Name: "valid complex type - deeply nested map of object with list of object with map(object) (PascalCase)",
			Content: `
variable "ComplexConfig" {
  type = map(object({
    ConfigName = string
    Rules = list(object({
      RuleType = string
      Targets = map(object({
	TargetId = string
	Metadata = object({
	  CreatedAt = string
	  OwnerId   = string
	})
      }))
    }))
  }))
  description = "valid."
}
`,
			Config:   testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{},
		},
		{
			Name: "invalid complex type - deeply nested map of object with list of object with map(object) (non-PascalCase)",
			Content: `
		variable "complex_config" {
		  type = map(object({
		    config_name = string
		    rules = list(object({
		      rule_type = string
		      targets = map(object({
			target_id = string
			metadata = object({
			  created_at = string
			  owner_id   = string
			})
		      }))
		    }))
		  }))
		  description = "invalid."
		}
		`,
			Config: testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "variable `complex_config` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 2, Column: 3},
						End:      hcl.Pos{Line: 2, Column: 28},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.config_name` - attribute `config_name` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.rule_type` - attribute `rule_type` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.targets.metadata.created_at` - attribute `created_at` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.targets.metadata.owner_id` - attribute `owner_id` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.targets.metadata` - attribute `metadata` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.targets.target_id` - attribute `target_id` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules.targets` - attribute `targets` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
					},
				},
				{
					Rule:    rule,
					Message: "variable `complex_config` path `complex_config.rules` - attribute `rules` must match the following custom_format: PascalCase (Starts with uppercase, no underscores allowed)",
					Range: hcl.Range{
						Filename: "main.tf",
						Start:    hcl.Pos{Line: 3, Column: 5},
						End:      hcl.Pos{Line: 15, Column: 8},
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

const testTerraformVarsObjectKeysNamingConventions_snakeCase = `
rule "terraform_vars_object_keys_naming_conventions" {
  enabled = true
  format  = "snake_case"
}
`

const testTerraformVarsObjectKeysNamingConventions_mixedSnakeCase = `
rule "terraform_vars_object_keys_naming_conventions" {
  enabled = true
  format  = "mixed_snake_case"
}
`

const testTerraformVarsObjectKeysNamingConventions_customFormat_pascalCase = `
rule "terraform_vars_object_keys_naming_conventions" {
  enabled = true
  custom_format_key = "PascalCase"

  custom_formats = {
    PascalCase = {
      regex       = "^[A-Z][a-zA-Z0-9]*$"
      description = "PascalCase (Starts with uppercase, no underscores allowed)"
    }
  }
}
`
