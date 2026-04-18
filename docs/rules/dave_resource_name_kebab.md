# dave_resource_name_kebab

Ensures `name` and `name_prefix` attributes use kebab-case (lowercase letters, numbers, and dashes only).

## Why

Cloud resource names appear in ARNs, URLs, and CLI output. Kebab-case is the most common convention for these external-facing identifiers and avoids issues with case-insensitive systems.

## Examples

```hcl
# ❌ Invalid
resource "aws_s3_bucket" "main" {
  name = "my_application_storage"
}

resource "aws_iam_role" "main" {
  name = "AdminRole"
}

# ✅ Valid
resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}

resource "aws_iam_role" "main" {
  name_prefix = "admin-role-"
}
```
