# dave_s3_no_inline_config

Flags deprecated inline configuration blocks and attributes on `aws_s3_bucket`.

## Why

Since AWS provider v4 (February 2022), S3 bucket sub-resources were split into dedicated resources. The inline properties on `aws_s3_bucket` are deprecated and will be removed in a future major release. Using the separate resources now avoids a forced rewrite later.

## Flagged blocks

`versioning`, `logging`, `server_side_encryption_configuration`, `lifecycle_rule`, `cors_rule`, `website`, `replication_configuration`, `object_lock_configuration`, `grant`

## Flagged attributes

`acl`, `policy`, `acceleration_status`, `request_payer`

## Examples

```hcl
# ❌ Invalid — inline versioning block
resource "aws_s3_bucket" "this" {
  bucket = "my-bucket"

  versioning {
    enabled = true
  }
}

# ✅ Valid — separate resource
resource "aws_s3_bucket" "this" {
  bucket = "my-bucket"
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.id

  versioning_configuration {
    status = "Enabled"
  }
}
```
