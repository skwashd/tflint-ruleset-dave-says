# dave_cloudwatch_log_retention

Ensures CloudWatch log groups have `retention_in_days` set to the expected value.

**Fixable:** Yes (when the attribute exists but has the wrong value)

## Why

Unbounded log retention wastes money. Log groups without explicit retention keep logs forever. 30 days is a sensible default — override only when there is a documented compliance reason for a different value.

## Configuration

The expected retention is configurable. Default is 30 days.

```hcl
rule "dave_cloudwatch_log_retention" {
  enabled        = true
  retention_days = 14
}
```

## Autofix

When run with `tflint --fix`, this rule will automatically correct `retention_in_days` values that don't match the configured expectation. If the attribute is missing entirely, the rule emits an issue but does not autofix (to avoid guessing indentation).

```
$ tflint --fix
1 issue(s) found:

Warning: [Fixed] CloudWatch log group retention_in_days is 90, expected 30. (dave_cloudwatch_log_retention)
  on cloudwatch.tf line 3:
   3:   retention_in_days = 90
```

## Examples

```hcl
# ❌ Invalid — missing retention
resource "aws_cloudwatch_log_group" "this" {
  name = "/aws/lambda/my-function"
}

# ❌ Invalid — non-standard retention (auto-fixable)
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 90
}

# ✅ Valid
resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/lambda/my-function"
  retention_in_days = 30
}
```
