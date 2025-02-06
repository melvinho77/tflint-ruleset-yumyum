# terraform_meta_arguments

Check the sequences and format of 'source', 'count', 'for_each', 'providers' and
'provider' meta arguments in Terraform modules, resources and data sources.

## Terraform modules
### Format
- Check the beginning arguments sequences in Terraform modules as the following:
  1. module definition *(end with newline)*
  2. source *(end with newline)*
  3. -- if count or for_each exist --
      1. count/for_each *(end with newline)*
      2. *(extra newline)*
  4. -- if providers exist --
      1. providers *(end with newline)*
      2. *(extra newline)*
  5. other attributes/blocks

### Valid example
```hcl
module "alicloud_ecs_instances" {
  source = "./alicloud-ecs-instance/"

  count = 3

  providers = {
    alicloud = alicloud.ecs
  }

  # ...
}
```

```hcl
module "aws_ec2_instance" {
  source = "./aws-ec2-instance/"

  providers = {
    aws = aws.ec2
  }

  # ...
}
```

## Terraform resources and data sources
### Format
- Check the beginning arguments sequences in Terraform resources and data sources
  as the following:
  1. resource definition *(end with newline)*
  2. -- if count or for_each exist --
      1. count/for_each *(end with newline)*
      2. *(extra newline)*
  3. -- if provider exist --
      1. provider *(end with newline)*
      2. *(extra newline)*
  4. other attributes/blocks
- `lifecycle{}` block must be placed as last block at the end of the resource without extra new lines.

## Valid example
```hcl
resource "aws_ec2_instance" "my_instance" {
  count = 3

  provider = aws.ec2

  # ...

  lifecycle {
    create_before_destroy = true
  }
}
```

```hcl
resource "aws_ec2_instance" "my_instance" {
  # ...

  lifecycle {
    create_before_destroy = true
  }
}
```

```hcl
data "aws_ec2_instance" "my_instance" {
  count = 3

  provider = aws.ec2

  # ...

  lifecycle {
    create_before_destroy = true
  }
}
```
