# dave_variable_region

Flags any variable named `region`.

## Why

Passing region as a variable creates a risk of mismatch with the provider's configured region. Use the provider's default region and `data.aws_region.current` when you need to reference it.

## Examples

```hcl
# ❌ Invalid
variable "region" {
  type    = string
  default = "us-east-1"
}

# ✅ Valid — use the provider's region
data "aws_region" "current" {}
# Reference: data.aws_region.current.name
```
