# dave_label_snake

Ensures labels on resource, data, ephemeral, module, output, and variable blocks use snake_case.

## Why

Terraform convention is snake_case for identifiers. Mixing camelCase or kebab-case breaks consistency and makes references harder to type correctly.

## Examples

```hcl
# ❌ Invalid
resource "aws_s3_bucket" "myStorage" {}
resource "aws_s3_bucket" "user-data" {}

# ✅ Valid
resource "aws_s3_bucket" "my_storage" {}
resource "aws_s3_bucket" "user_data" {}
variable "storage_name" {}
output "bucket_arn" {}
```
