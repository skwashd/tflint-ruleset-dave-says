# dave_output_must_be_in_outputs_file

Ensures all `output` blocks are declared in `outputs.tf`.

## Why

Consistent file organization makes modules predictable. When outputs are scattered across multiple files, consumers have to search for them. Keeping all outputs in `outputs.tf` mirrors the convention of keeping variables in `variables.tf`.

## Examples

```hcl
# ❌ Invalid — output in main.tf
# main.tf
output "bucket_arn" {
  value = aws_s3_bucket.this.arn
}

# ✅ Valid — output in outputs.tf
# outputs.tf
output "bucket_arn" {
  description = "The ARN of the S3 bucket."
  value       = aws_s3_bucket.this.arn
}
```
