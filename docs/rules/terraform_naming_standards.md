# terraform_naming_standards

Enforces naming conventions specifically for `variable` type blocks, validating both variable names and any nested object field names based on the configured format (e.g., `snake_case` or custom regex). This rule builds on the existing [`terraform_naming_convention`](https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/rules/terraform_naming_convention.md) from `tflint-ruleset-terraform` by extending its coverage to variable's sub-attributes, and not just the variable's name.

## Configuration

| Name           | Default      | Value                                                                                                     |
| -------------- | ------------ | --------------------------------------------------------------------------------------------------------- |
| enabled        | true         | Bool                                                                                                      |
| format         | `snake_case` | `snake_case`, `mixed_snake_case`, `none`, or a custom format defined using the `custom_formats` attribute |
| custom         | ""           | String representation of a golang regular expression that the block name must match                       |
| custom_formats | {}           | Definition of custom formats that can be used in the `format` attribute                                   |
| variable       | {}           | Block settings to override naming conventions for variables                                               |

#### `format`

The `format` option defines the allowed formats for the block label. This option accepts one of the following values:

- `snake_case` - standard snake_case format - all characters must be lower-case, and underscores are allowed.
- `mixed_snake_case` - modified snake_case format - characters may be upper or lower case, and underscores are allowed.
- `none` - if this option is selected, it does not perform any regex checking on the `variable` blocks.

#### `custom`

- The `custom` option defines a custom regex that the identifier must match. The regex must follow the Golang regex syntax. It is to allow more fine-grained control over predefined identifier patterns and structure.

#### `custom_formats`

- The `custom_formats` option defines additional formats that can be used in the `format` option. Like `custom`, it allows you to define a custom regular expression that the identifier must match, but it also lets you define a description that will be shown when the check fails. Additionally, the use of custom regex is also allowed.

- This attribute is a map, where the keys are the identifiers of the custom formats, and the values are objects with a regex and a description key.

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

### Default - enforce `snake_case` for all blocks

#### Rule configuration

```hcl
rule "terraform_naming_standards" {
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

Warning: variable name `invalidName` must match the following format: snake_case (terraform_naming_standards)

  on main.tf line 1:
   1: variable "invalidName" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_naming_standards.md

Warning: variable name `invalid_object.fooBar` must match the following format: snake_case (terraform_naming_standards)

  on main.tf line 6:
   6:   type = object({
   7:     foo_bar = string
   8:     fooBar  = bool
   9:   })

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_naming_standards.md
```

### Custom regex for all blocks

#### Rule configuration

```hcl
rule "terraform_naming_standards" {
  enabled = true

  custom = "^[a-zA-Z]+([_-][a-zA-Z]+)*$"
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

Warning: variable name `Invalid_Name_With_Number123` must match the following RegExp: ^[a-zA-Z]+([_-][a-zA-Z]+)*$ (terraform_naming_standards)

  on main.tf line 1:
   1: variable "Invalid_Name_With_Number123" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_naming_standards.md
```

### Custom format for all blocks

#### Rule configuration

```hcl
rule "terraform_naming_standards" {
  enabled = true
  format  = "custom_format"

  custom_formats = {
    custom_format = {
      description = "Custom Format"
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

Warning: variable name `Invalid_Name_With_Number123` must match the following format: Custom Format (terraform_naming_standards)

  on main.tf line 1:
   1: variable "Invalid_Name_With_Number123" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_naming_standards.md
```
