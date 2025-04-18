# terraform_required_tags

Check whether `resources` have required tags labeled for resources with the `tag` block defined. Additionally, for AWS resources, the `Name` tag key must be included for any resources with the `tag` block defined.

## Configuration

| Name               | Default | Value          |
| ------------------ | ------- | -------------- |
| enabled            | true    | Bool           |
| tags               | []      | List of string |
| excluded_resources | []      | List of string |

#### `tags`

The `tags` option defines the list of tags that needs to be included for any resources that has the `tags` block.

#### `excluded_resources`

The `excluded_resources` option defines the list of resources type to be ignored in ths rule checking.

## Example

### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled            = true
  tags               = ["example_tag1", "example_tag2", "example_tag3"]
}
```

#### Sample terraform source file

```hcl
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

### Disable for specificed resources

#### Rule configuration

```hcl
rule "terraform_required_tags" {
  enabled            = true
  tags               = ["example_tag1", "example_tag2", "example_tag3"]
  excluded_resources = ["my_excluded_resource"]
}
```

#### Sample terraform source file

```hcl
// resource "my_excluded_resource" will not be enforced

resource "my_excluded_resource" "my_resource_name" {
  name = "test"

  tags = {
    example_tag1 = "value1"
    example_tag2 = "value2"
  }
}
```
