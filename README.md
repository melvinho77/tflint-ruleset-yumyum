# TFLint Ruleset for ST Terraform Standrds

TFLint ruleset plugin for Terraform Language

This ruleset focus on the best practice standards about ST Terraform Language practice.

## Requirements

- TFLint v0.56.0+
- Go v1.24

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

Once installed, you can use the plugin in your Terraform modules by creating a `.tflint.hcl` file that contains the following content:

```hcl
plugin "yumyum" {
  enabled = true
}
```

## Rules

| Rule                          | Description                                                                                                                                                                                                                                                                                                                         |
| ----------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| terraform_any_type_variables  | Disallow `variable` declarations with type `any`                                                                                                                                                                                                                                                                                    |
| terraform_meta_arguments      | Ensure correct ordering and formatting of `source`, `count`, `for_each`, `providers`, and `provider` in `module`, `resource`, and `data` blocks.                                                                                                                                                                                    |
| terraform_module_dependencies | Ensure `module` sources are pinned to a specific version using `?ref=` or `?rev=` in source URLs.                                                                                                                                                                                                                                   |
| terraform_naming_standards    | Extends [`terraform_naming_convention`](https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/rules/terraform_naming_convention.md) by enforcing naming conventions not just for `variable` names, but also for nested object field names, based on a configured format like `snake_case` or a custom regex. |
| terraform_required_tags       | Checks if resources include required tags in their `tags` block. For AWS, enforces presence of the `Name` tag as well.                                                                                                                                                                                                              |
| terraform_required_variables  | Ensures all variables listed in `required_vars` are declared in the Terraform module.                                                                                                                                                                                                                                               |
|                               |
