# dave_iam_no_inline_policy

Flags use of `aws_iam_role_policy`, `aws_iam_user_policy`, and `aws_iam_group_policy` resources.

## Why

These resources create inline policies embedded directly on the IAM identity. Inline policies are not reusable across roles and are harder to audit in the IAM console. Use `aws_iam_policy` with the corresponding `_policy_attachment` resource instead.

## Examples

```hcl
# ❌ Invalid — inline policy
resource "aws_iam_role_policy" "lambda" {
  name   = "app-prod-lambda"
  role   = aws_iam_role.lambda.id
  policy = data.aws_iam_policy_document.lambda.json
}

# ✅ Valid — managed policy with attachment
resource "aws_iam_policy" "lambda" {
  name   = aws_iam_role.lambda.name
  policy = data.aws_iam_policy_document.lambda.json
}

resource "aws_iam_role_policy_attachment" "lambda" {
  role       = aws_iam_role.lambda.name
  policy_arn = aws_iam_policy.lambda.arn
}
```
