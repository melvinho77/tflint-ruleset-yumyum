# terraform_vars_object_keys_naming_conventions

Enforces naming conventions specifically for `variable` type blocks, validating both variable names and any nested object field names based on the configured format (e.g., `snake_case` or custom regex). This rule builds on the existing [`terraform_naming_convention`](https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/rules/terraform_naming_convention.md) from `tflint-ruleset-terraform` by extending its coverage to variable's sub-attributes, and not just the variable's name.

## Configuration

| Name              | Default      | Value                                                                                                      |
| ----------------- | ------------ | ---------------------------------------------------------------------------------------------------------- |
| enabled           | true         | `true` or `false` - Enable or disables the rule.                                                           |
| format            | `snake_case` | `snake_case`, `mixed_snake_case`, `none`.                                                                  |
| custom_format_key | ""           | The key from `custom_formats` to use for custom regex matching (e.g., `PascalCase`, `camelCase`)           |
| custom_formats    | {}           | A map of custom formats, where each key defines a format with `regex` (string) and `description` (string). |

#### `format`

The `format` option defines the allowed predefined formats for the tflint rule config. This option accepts one of the following values:

- `snake_case` - standard snake_case format - all characters must be lower-case, and underscores are allowed.
- `mixed_snake_case` - modified snake_case format - characters may be upper or lower case, and underscores are allowed.
- `none` - if this option is selected, it does not perform any regex checking on the `variable` blocks.

#### `custom_format_key`

- This option selects a custom format from `custom_formats`. The selected format will be applied for validation using its defined regex pattern.
- For example, to use and apply a custom format:

```hcl
rule "terraform_naming_standards" {
  enabled           = true
  custom_format_key = "upper_snake"
  custom_formats    = {
    upper_snake = {
      regex       = "xxx"
      description = "xxx"
    }
  }
}
```

#### `custom_formats`

- Defines one or more custom naming conventions using regex and a description to guide users. Useful for enforcing team-specific or organization-specific naming patterns.

- Each entry in this map should have:

  - `regex` - A Golang-compatible regular expression string.
  - `description` - An explanation of the custom format pattern.

- Example `custom_formats` definition:

```hcl
  custom_formats = {
    upper_snake = {
      regex       = "^[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$"
      description = "UPPER_SNAKE_CASE format"
    }

    kebab = {
      regex       = "^[a-z0-9]+(-[a-z0-9]+)*$"
      description = "kebab-case format"
    }
  }
```

## Examples

### Default - enforce `snake_case` as the default rule

#### Rule configuration

```hcl
rule "terraform_vars_object_keys_naming_conventions" {
  enabled = true
}
```

#### Sample terraform source file

```hcl
variable "invalidName" {
  type = string
}

variable "invalid_object" {
  type = object({
    foo_bar = string
    fooBar  = bool
  })
}

variable "valid_name" {
  type = string
}
```

```
$ tflint
2 issue(s) found:

Warning: variable `invalidName` must match the following predefined_format: snake_case (terraform_vars_object_keys_naming_conventions)

  on main.tf line 1:
   1: variable "invalidName" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md

Warning: variable `invalid_object` path `invalid_object.fooBar` - attribute `fooBar` must match the following predefined_format: snake_case (terraform_vars_object_keys_naming_conventions)

  on main.tf line 6:
   6:   type = object({
   7:     foo_bar = string
   8:     fooBar  = bool
   9:   })

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md
```

### Enforce predefined format rule - `mixed_snake_case`

#### Rule configuration

```hcl
rule "terraform_vars_object_keys_naming_conventions" {
  enabled = true
  format  = "mixed_snake_case"
}
```

#### Sample terraform source file

```hcl
variable "Invalid_Name_With_Multiple__Underscores" {
  type = string
}

variable "Name-With_Dash" {
  type = string
}
```

```
$ tflint
2 issue(s) found:

Warning: variable `Invalid_Name_With_Multiple__Underscores` must match the following predefined_format: mixed_snake_case (terraform_vars_object_keys_naming_conventions)

  on main.tf line 1:
   1: variable "Invalid_Name_With_Multiple__Underscores" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md

Warning: variable `Name-With_Dash` must match the following predefined_format: mixed_snake_case (terraform_vars_object_keys_naming_conventions)

  on main.tf line 5:
   5: variable "Name-With_Dash" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md
```

### Enforce a custom format

#### Rule configuration

```hcl
rule "terraform_vars_object_keys_naming_conventions" {
  enabled           = true
  custom_format_key = "custom_format"

  custom_formats = {
    custom_format = {
      description = "Custom Format [Alphabetic words separated by hyphens or underscores (e.g., 'my_variable', 'My-Variable')]"
      regex       = "^[a-zA-Z]+([_-][a-zA-Z]+)*$"
    }
  }
}
```

#### Sample terraform source file

```hcl
variable "Invalid_Name_With_Number123" {
  type = string
}

variable "Name-With_Dash" {
  type = string
}
```

```
$ tflint
1 issue(s) found:

Warning: variable `Invalid_Name_With_Number123` must match the following custom_format: Custom Format [Alphabetic words separated by hyphens or underscores (e.g., 'my_variable', 'My-Variable')] (terraform_vars_object_keys_naming_conventions)

  on main.tf line 1:
   1: variable "Invalid_Name_With_Number123" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/teraform_vars_object_keys_naming_conventions.md
```
