# dave_variable_has_description

Ensures all variables have a `description` attribute.

## Why

Without a description, anyone consuming the module must read the implementation to understand what a variable controls. This doesn't scale, especially in modules with many variables.

## Examples

```hcl
# ❌ Invalid — missing description
variable "bucket_name" {
  type = string
}

# ✅ Valid
variable "bucket_name" {
  description = "The prefix for the S3 bucket name."
  type        = string
}
```
