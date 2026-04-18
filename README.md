# TFLint Ruleset: Dave Says

A TFLint plugin that enforces custom Terraform coding standards for consistent naming conventions, proper code organization, and best practices.

## Rules

### Label Rules
* **[dave_label_min_length](docs/rules/dave_label_min_length.md):** Enforce minimum length for labels
* **[dave_label_no_type_substring](docs/rules/dave_label_no_type_substring.md):** Avoid redundant information in label values
* **[dave_label_snake](docs/rules/dave_label_snake.md):** Use snake_case for all labels

### Resource Name Rules
* **[dave_resource_name_kebab](docs/rules/dave_resource_name_kebab.md):** Use kebab-case for all names
* **[dave_resource_name_no_type_substring](docs/rules/dave_resource_name_no_type_substring.md):** Avoid redundant information in resource names

### Variable Rules
* **[dave_no_vpc_id_variable](docs/rules/dave_no_vpc_id_variable.md):** Do not use `vpc_id` as a variable name
* **[dave_variable_alphabetical_order](docs/rules/dave_variable_alphabetical_order.md):** Sort variables alphabetically
* **[dave_variable_has_description](docs/rules/dave_variable_has_description.md):** All variables must have a `description`
* **[dave_variable_has_type](docs/rules/dave_variable_has_type.md):** All variables must have a `type`
* **[dave_variable_must_be_in_variables_file](docs/rules/dave_variable_must_be_in_variables_file.md):** Only allow variables in `variables.tf`
* **[dave_variable_region](docs/rules/dave_variable_region.md):** Do not use `region` as a variable name

### File Organization Rules
* **[dave_output_must_be_in_outputs_file](docs/rules/dave_output_must_be_in_outputs_file.md):** Only allow outputs in `outputs.tf`

### AWS IAM Rules
* **[dave_aws_policy_no_jsonencode](docs/rules/dave_aws_policy_no_jsonencode.md):** Ensure policies use `aws_iam_policy_document` data sources
* **[dave_iam_no_inline_policy](docs/rules/dave_iam_no_inline_policy.md):** Use managed policies instead of inline policies

### AWS S3 Rules
* **[dave_s3_bucket_namespace](docs/rules/dave_s3_bucket_namespace.md):** S3 buckets must use account-regional namespace
* **[dave_s3_no_inline_config](docs/rules/dave_s3_no_inline_config.md):** Do not use deprecated inline S3 configuration blocks
* **[dave_s3_no_public_acl](docs/rules/dave_s3_no_public_acl.md):** S3 bucket ACLs must not allow public access

### AWS CloudWatch Rules
* **[dave_cloudwatch_log_retention](docs/rules/dave_cloudwatch_log_retention.md):** CloudWatch log groups must set `retention_in_days`

### AWS VPC Rules
* **[dave_security_group_no_inline_rules](docs/rules/dave_security_group_no_inline_rules.md):** Do not use inline `ingress`/`egress` blocks on security groups

## Installation

### Building from Source

```bash
git clone https://github.com/skwashd/tflint-ruleset-dave-says.git
cd tflint-ruleset-dave-says
go build -o tflint-ruleset-dave-says
```

### Installing the Plugin

```bash
mkdir -p ~/.tflint.d/plugins
cp tflint-ruleset-dave-says ~/.tflint.d/plugins/
```

## Configuration

Add to your `.tflint.hcl`:

```hcl
plugin "dave-says" {
  enabled = true
  version = "0.2.0"
  source  = "github.com/skwashd/tflint-ruleset-dave-says"
}
```

### Disabling Individual Rules

```hcl
rule "dave_cloudwatch_log_retention" {
  enabled = false
}

rule "dave_s3_bucket_namespace" {
  enabled = false
}
```

### Rule Configuration

Some rules accept per-rule configuration:

```hcl
rule "dave_cloudwatch_log_retention" {
  enabled        = true
  retention_days = 14  # default: 30
}
```

| Rule | Option | Type | Default | Description |
|------|--------|------|---------|-------------|
| `dave_cloudwatch_log_retention` | `retention_days` | int | `30` | Expected `retention_in_days` value |

## Autofix

Some rules support `tflint --fix` to automatically correct issues:

| Rule | Fix action |
|------|-----------|
| `dave_cloudwatch_log_retention` | Replaces wrong `retention_in_days` value with the configured expected value |
| `dave_s3_no_public_acl` | Replaces public ACL (`public-read`, `public-read-write`, `authenticated-read`) with `"private"` |

Autofix only applies when the attribute exists but has the wrong value. Missing attributes are flagged but not auto-fixed to avoid guessing indentation and placement.

Requires TFLint v0.47+.

## Provider Version Requirements

Some rules target features that require minimum provider versions:

| Rule | Minimum AWS Provider |
|------|---------------------|
| `dave_s3_bucket_namespace` | >= 6.37.0 |
| `dave_s3_no_inline_config` | >= 4.0.0 (deprecated since) |
| `dave_security_group_no_inline_rules` | >= 5.0.0 |

## Development

### Requirements

- Go 1.23 or newer
- TFLint v0.59.0 or newer

### Building

```bash
go build -o tflint-ruleset-dave-says
```

### Testing

```bash
go test ./rules/... -v
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
