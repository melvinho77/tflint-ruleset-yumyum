# terraform_any_type_variables

Disallow `variable` declarations with type `any`.

## Configuration

| Name        | Default | Value          |
| ----------- | ------- | -------------- |
| enabled     | true    | Boolean        |
| ignore_vars | []      | List of string |

#### `ignore_vars`

The `ignore_vars` option defines the list of variables name to be ignored in this
rule checking.

## Example

### Default - enforce disallow `variable` declarations with type `any`.

#### Rule configuration

```hcl
rule "terraform_any_type_variables" {
  enabled = true
}
```

#### Sample terraform source file

```hcl
variable "my_var" {
  type = any
}
```

```
$ tflint
1 issue(s) found:

Warning: variable 'my_var' has 'any' type declared (terraform_any_type_variables)

  on variables.tf line 2:
 2:         type = any

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_any_type_variables.md
```

### Disable for specified variables

#### Rule configuration

```hcl
rule "terraform_any_type_variables" {
  enabled = true

  ignore_vars = ["my_var"]
}
```

#### Sample terraform source file

```hcl
// variable 'my_var' will not be enforced
variable "my_var" {
  type = any
}
```
