# dave_resource_name_no_type_substring

Prevents `name` and `name_prefix` attributes from containing words that appear in the resource type.

## Why

The resource type already communicates what the thing is. Repeating it in the name wastes characters and adds no clarity. `name = "admin-policy"` on an `aws_iam_role` tells you it's a role for admin policy — but `name = "admin-role-policy"` just repeats "role".

## Examples

```hcl
# ❌ Invalid — "role" appears in both the type and the name
resource "aws_iam_role" "main" {
  name = "admin-role-policy"
}

# ❌ Invalid — "bucket" appears in both
resource "aws_s3_bucket" "main" {
  name = "my-application-bucket"
}

# ✅ Valid
resource "aws_iam_role" "main" {
  name = "admin-policy"
}

resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}
```
