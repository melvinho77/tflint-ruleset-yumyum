# terraform_module_source_version

Check whether `module` sources have explicitly pinned to a semantic versioning using the `?ref=` or `?rev=` query parameters in their source URLs. The `allowed_versions` can be specified as a list of regex patterns to permit flexible versioning schemes beyond strict semantic versions.

## Configuration

| Name             | Default | Value          |
| ---------------- | ------- | -------------- |
| enabled          | `true`  | Bool           |
| allowed_versions | `[]`    | List of string |

#### `allowed_version`

The `allowed_version` option defines a list of regular expressions / exact strings used to validate the `?ref=` or `?rev=` versioning in the source URL. Only versions matching on of these expressions will be allowed. Example Go regular expressions are as follow:

- `^bugfix/\\d+$`
- `^feature/\\d+$`

## Example

### Rule configuration

```hcl
rule "terraform_module_source_version" {
  enabled          = true
  allowed_versions = ["^bugfix/\\d+$", "^feature/\\d+$", "bugfix/test"]
}
```

#### Sample terraform source file

```hcl
module "my_module_1" {
  source = "git://gitlab.example.com/test.git?ref=main"
}

module "my_module_2" {
  source = "git://gitlab.example.com/test.git?ref=bugfix/test"
}

// The following modules demonstrate valid pinned references:
// - local path reference
// - git references pinned by semver tag (v1.2.3)
// - git references pinned by allowed non-semver branch names (bugfix/1234, feature/1234)

module "my_module_3" {
  source = "../test"
}

module "my_module_4" {
  source = "git://gitlab.example.com/test.git?ref=v1.2.3"
}

module "my_module_5" {
  source = "git://gitlab.example.com/test.git?ref=bugfix/1234"
}

module "my_module_6" {
  source = "git://gitlab.example.com/test.git?ref=feature/1234"
}

module "my_module_7" {
  source = "git://gitlab.example.com/test.git?ref=bugfix/test"
}
```

```
2 issue(s) found:

Warning: module 'my_module_1' source 'git://gitlab.example.com/test.git?ref=main' [ref='main'] does not match any allowed_versions pattern (terraform_module_source_version)

  on main.tf line 2:
   2:   source = "git://gitlab.example.com/test.git?ref=main"

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_source_version.md

Warning: module 'my_module_2' source 'git://gitlab.example.com/test.git?ref=bugfix/test' [ref='bugfix/test'] does not match any allowed_versions pattern (terraform_module_source_version)

  on main.tf line 6:
   6:   source = "git://gitlab.example.com/test.git?ref=bugfix/test"

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_module_source_version.md
```
