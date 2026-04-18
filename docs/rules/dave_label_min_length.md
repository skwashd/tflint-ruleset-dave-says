# dave_label_min_length

Ensures resource, data, and ephemeral block labels are at least 3 characters long.

## Why

Single and two-character labels are rarely descriptive enough to convey purpose. Longer names improve readability and make `grep`/search more effective.

The well-known abbreviations `db` and `s3` are allowed as exceptions.

## Examples

```hcl
# ❌ Invalid — too short
resource "aws_s3_bucket" "ab" {}
resource "aws_s3_bucket" "a" {}

# ✅ Valid
resource "aws_s3_bucket" "assets" {}
data "aws_iam_policy_document" "s3" {}
resource "aws_db_instance" "db" {}
```
