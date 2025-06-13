# terraform_required_tags

This rule checks that all Terraform resources with a `tags` block include the required tag keys defined in the rule configuration. It supports both direct tag maps and `merge()` expressions, specifically allowing `merge(local.tags, {...})`, and will evaluate and combine all tag keys before validating them. If `local.tags` is missing, it reports an issue. Additionally, for AWS resources, it enforces the presence of a `Name` tag. Unsupported expressions or function calls in tags will trigger a warning. Resources listed in the excluded list are skipped.

## Configuration

| Name               | Default                                                                                           | Value          |
| ------------------ | ------------------------------------------------------------------------------------------------- | -------------- |
| enabled            | true                                                                                              | Bool           |
| tags               | ["brand", "env", "project", "devops_project_kind", "devops_project_group", "devops_project_name"] | List of string |
| excluded_resources | []                                                                                                | List of string |

#### `tags`

The `tags` option defines the list of tags that needs to be included for any resources that has the `tags` block. Defaults to the following list:

```hcl
tags = [
  "brand",
  "env",
  "project",
  "devops_project_kind",
  "devops_project_group",
  "devops_project_name",
]
```

#### `excluded_resources`

The `excluded_resources` option defines the list of resources type to be ignored in ths rule checking. Defaults to an empty list.

## Example

### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled            = true
  tags               = ["example_tag1", "example_tag2", "example_tag3"]
}
```

### Sample terraform source file

```hcl
local {
  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
    example_tag3 = "value3"
  }
}

resource "my_resource" "my_resource_name" {
  name = "test"

  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: my_resource 'my_resource_name' is missing required tags: [example_tag3] (terraform_required_tags)

  on main.tf line 4:
   4:   tags = {
   5:     example_tag1 = "value1"
   6:     example_tag2 = "value2"
   7:   }

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_tags.md
```

## Disable for specificed resources

### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled            = true
  tags               = ["example_tag1", "example_tag2", "example_tag3"]
  excluded_resources = ["my_excluded_resource"]
}
```

### Sample terraform source file

```hcl
local {
  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
    example_tag3 = "value3"
  }
}

// resource "my_excluded_resource" will not be enforced
resource "my_excluded_resource" "my_resource_name" {
  name = "test"

  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
  }
}
```

## Direct use of local `tags` variable

### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled = true
  tags               = ["example_tag1", "example_tag2", "example_tag3"]
}
```

### Sample terraform source file

```hcl
locals {
  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
  }
}

resource "my_excluded_resource" "my_resource_name" {
  name = "test"

  tags = local.tags
}
```

```hcl
$ tflint
1 issue(s) found:

Warning: my_excluded_resource 'my_resource_name' is missing required tags: [example_tag3] (terraform_required_tags)

  on main.tf line 11:
  11:   tags = local.tags

Reference: https://github.com/styumyum/tflint-ruleset-yumyum/docs/rules/terraform_required_tags.md
```

## Usage of function call `merge` with locals `tags`

### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled = true
  tags    = ["example_tag1", "example_tag2", "example_tag3"]
}
```

### Sample terraform source file

```hcl
locals {
  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
  }
}

resource "my_excluded_resource" "my_resource_name" {
  name = "test"

  tags = merge(local.tags, {
    example_tag3 = "value3"
  })
}
```
