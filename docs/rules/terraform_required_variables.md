# terraform_required_variables

Check whether the list of variables declared in `required_vars` are also declared in the Terraform module.

## Configuration

| Name          | Default                                       | Value          |
| ------------- | --------------------------------------------- | -------------- |
| enabled       | true                                          | Bool           |
| required_vars | ["cloud_creds", "module_info", "module_tmpl"] | List of string |

### `required_vars`

The `required_vars` option defines the list of variables that is mandatory to be defined in the terraform module.

## Example

#### Rule configuration

```hcl
rule "terraform_required_variables" {
  enabled = true
  required_vars = ["var1", "var2", "var3"]
}
```

#### Sample terraform source file

```hcl
variable "cloud_creds" {
  type = string
}
```

```
$ tflint
2 issue(s) found:

Warning: required variable(s) not declared: var1, var2, var3, module_info, module_tmpl (terraform_required_variables)

  on  line 1:
   (source code not available)

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_variables.md

Warning: variable `cloud_creds` is missing the `sensitive` attribute (terraform_required_variables)

  on variables.tf line 1:
   1: variable "cloud_creds" {

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_variables.md
```

## Selected variables that is mandatory

#### Rule configuration

```
rule "terraform_required_variables" {
  enabled = true
  required_vars = ["var1", "var2"]
}
```

#### Sample terraform source file

```hcl
variable "var1" {
  type = string
}
```

```
$ tflint
1 issue(s) found:

Warning: required variable(s) not declared: var2, cloud_creds, module_info, module_tmpl (terraform_required_variables)

  on  line 1:
   (source code not available)

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_variables.md
```
