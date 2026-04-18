# dave_variable_must_be_in_variables_file

Ensures all variable blocks are declared in `variables.tf`.

## Why

Scattering variables across files forces readers to search the entire module to find inputs. A single `variables.tf` file makes the module's interface immediately visible.

## Examples

```hcl
# ❌ Invalid — variable in main.tf
# main.tf
variable "name" {
  type = string
}

# ✅ Valid — variable in variables.tf
# variables.tf
variable "name" {
  type = string
}
```
