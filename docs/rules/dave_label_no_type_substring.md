# dave_label_no_type_substring

Prevents labels from containing words that already appear in the resource type.

## Why

The resource type is always visible alongside the label. Repeating words from the type adds noise without information. `aws_s3_bucket.user_data` is clearer than `aws_s3_bucket.user_bucket`.

## Examples

```hcl
# ❌ Invalid — label repeats words from type
resource "aws_s3_bucket" "user_bucket" {}
resource "aws_iam_role" "admin_role" {}
resource "aws_s3_bucket" "s3_storage" {}

# ✅ Valid
resource "aws_s3_bucket" "user_data" {}
resource "aws_iam_role" "application" {}
```
