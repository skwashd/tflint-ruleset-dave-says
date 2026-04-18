# dave_aws_policy_no_jsonencode

Flags use of `jsonencode()` in policy attributes on AWS resources.

## Why

Using `jsonencode()` for IAM policies bypasses Terraform's ability to validate policy structure at plan time. The `aws_iam_policy_document` data source provides type safety, automatic JSON formatting, and better diff output.

## Examples

```hcl
# ❌ Invalid — jsonencode in policy
resource "aws_iam_role" "main" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Principal = { Service = "lambda.amazonaws.com" }
    }]
  })
}

# ✅ Valid — policy document data source
resource "aws_iam_role" "main" {
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }
  }
}
```

Non-AWS resources are not affected by this rule.
