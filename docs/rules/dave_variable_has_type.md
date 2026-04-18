# dave_variable_has_type

Ensures all variables have an explicit `type` constraint.

## Why

Without a type constraint, Terraform cannot catch invalid inputs at plan time, pushing errors to apply where they are more expensive and disruptive.

## Examples

```hcl
# ❌ Invalid — missing type
variable "environment" {
  description = "The deployment environment."
}

# ✅ Valid
variable "environment" {
  description = "The deployment environment."
  type        = string
}
```
