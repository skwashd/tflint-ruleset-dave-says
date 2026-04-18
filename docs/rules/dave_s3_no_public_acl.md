# dave_s3_no_public_acl

Flags S3 bucket ACLs that allow public access.

**Fixable:** Yes (replaces public ACL with `"private"`)

## Why

Public S3 buckets are consistently one of the top causes of data breaches in AWS. Use CloudFront for web access instead — it adds caching, TLS, and access controls.

## Flagged values

`public-read`, `public-read-write`, `authenticated-read`

## Autofix

When run with `tflint --fix`, this rule replaces any public ACL with `"private"`.

```
$ tflint --fix
1 issue(s) found:

Error: [Fixed] S3 bucket ACL is "public-read". Public buckets are a leading cause of data breaches. Use CloudFront for web access. (dave_s3_no_public_acl)
  on s3.tf line 3:
   3:   acl = "public-read"
```

## Examples

```hcl
# ❌ Invalid (auto-fixed to "private")
resource "aws_s3_bucket_acl" "this" {
  bucket = aws_s3_bucket.this.id
  acl    = "public-read"
}

# ✅ Valid
resource "aws_s3_bucket_acl" "this" {
  bucket = aws_s3_bucket.this.id
  acl    = "private"
}
```
