# terraform_module_dependencies

Check whether `module` sources have explicitly pinned to a sepcific version using the `?ref=` or `?rev=` query parameters in their source URLs.

## Configuration

| Name             | Default | Value          |
| ---------------- | ------- | -------------- |
| enabled          | true    | Bool           |
| style            | ""      | String         |
| default_branches | []      | List of string |

#### `style`

The `style` option determines the enforcement policy for how module versions should be pinned:

- `semver` requires the version is pinned according to Semantic Versioning string (e.g. v1.0.0).
- `flexible` allows more flexibility, but disallow the usage of default branchse (like `main` or `master`) as the pin target. Any other revision (e.g., a feature branch or commit hash) is acceptable.

#### `default_branches`

The `default_branches` option defines a list of string which are not allowed to be used as the pinned version.

## Example

### Rule configuration

#### `style = "flexible"`

```hcl
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "flexible"
  default_branches = ["main", "master"]
}
```

#### Sample terraform source file

```hcl
module "my_module_1" {
  source = "git://gitlab.example.com/test.git?ref=main"

}

module "my_module_2" {
  // This will not enforce a warning because the style is set to `flexible`
  source = "git://gitlab.example.com/test.git?ref=feature/1234"
}
```

```
1 issue(s) found:

Warning: module 'my_module_1' source 'git://gitlab.example.com/test.git?ref=main' uses a default branch as ref (main) (terraform_module_dependencies)

  on main.tf line 2:
   2:   source = "git://gitlab.example.com/test.git?ref=main"

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_dependencies.md
```

#### `style = "semver"`

```hcl
rule "terraform_module_dependencies" {
  enabled          = true
  style            = "semver"
  default_branches = ["main", "master"]
}
```

#### Sample terraform source file

```hcl
module "my_module_1" {
  source = "git://gitlab.example.com/test.git?ref=main"

}

module "my_module_2" {
  source = "git://gitlab.example.com/test.git?ref=feature/1234"

}

module "my_module_3" {
  source = "git://gitlab.example.com/test.git?ref=v1.2.3"

}
```

```
$ tflint
2 issue(s) found:

Warning: module 'my_module_1' source 'git://gitlab.example.com/test.git?ref=main' uses a ref which is not a semantic version string (terraform_module_dependencies)

  on main.tf line 2:
   2:   source = "git://gitlab.example.com/test.git?ref=main"

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_dependencies.md

Warning: module 'my_module_2' source 'git://gitlab.example.com/test.git?ref=feature/1234' uses a ref which is not a semantic version string (terraform_module_dependencies)

  on main.tf line 7:
   7:   source = "git://gitlab.example.com/test.git?ref=feature/1234"

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_dependencies.md
```
