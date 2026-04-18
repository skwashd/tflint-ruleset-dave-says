# dave_s3_bucket_namespace

Ensures S3 buckets use account-regional namespace (`bucket_namespace = "account-regional"`).

## Why

S3 bucket names in the global namespace are first-come-first-served across all AWS accounts. This causes naming collisions, makes bucket names unpredictable across environments, and opens the door to confused-deputy attacks. Account regional namespaces scope bucket names to your account and region.

Requires AWS provider >= 6.37.0.

## Examples

```hcl
# ✅ Valid
resource "aws_s3_bucket" "this" {
  bucket           = format("%s-%s-%s-an", "my-app", data.aws_caller_identity.current.account_id, data.aws_region.current.region)
  bucket_namespace = "account-regional"
}

# ❌ Invalid — missing bucket_namespace
resource "aws_s3_bucket" "this" {
  bucket = "my-app-bucket"
}

# ❌ Invalid — wrong namespace value
resource "aws_s3_bucket" "this" {
  bucket           = "my-app-bucket"
  bucket_namespace = "global"
}
```
