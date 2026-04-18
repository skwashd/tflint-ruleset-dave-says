# dave_variable_alphabetical_order

Ensures variables within each file are sorted alphabetically by name.

## Why

Alphabetical ordering makes it easy to find a variable in a long list and produces predictable diffs — new variables slot into a known position rather than landing wherever the author happened to add them.

## Examples

```hcl
# ❌ Invalid — out of order
variable "environment" {
  type = string
}

variable "application_name" {
  type = string
}

# ✅ Valid
variable "application_name" {
  type = string
}

variable "environment" {
  type = string
}
```
