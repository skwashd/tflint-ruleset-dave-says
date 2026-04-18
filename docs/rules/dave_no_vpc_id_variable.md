# dave_no_vpc_id_variable

Flags any variable named `vpc_id`.

## Why

Passing both `vpc_id` and `subnet_ids` creates a class of bugs where they don't match — and the error surfaces deep in the apply, not at plan time. Derive VPC ID from the first subnet instead.

## Examples

```hcl
# ❌ Invalid
variable "vpc_id" {
  description = "The VPC ID"
  type        = string
}

# ✅ Valid — derive from subnets
data "aws_subnet" "first" {
  id = var.subnet_ids[0]
}
# Reference: data.aws_subnet.first.vpc_id
```
