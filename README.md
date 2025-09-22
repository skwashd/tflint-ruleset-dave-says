# TFLint Ruleset: Dave Says

A TFLint plugin that enforces custom Terraform coding standards for consistent naming conventions, proper code organization, and best practices.

## Rules

* **[dave_aws_policy_no_jsonencode](#dave_aws_policy_no_jsonencode):** Ensure policies use aws_iam_policy_documents
* **[dave_label_min_length](#dave_label_min_length):** Enforce minimum length for labels
* **[dave_label_no_type_substring](#dave_label_no_type_substring):** Avoid redundant information in label values
* **[dave_label_snake](#dave_label_snake):** Use snake_case for all labels
* **[dave_resource_name_kebab](#dave_resource_name_kebab):** Use kebab-case for all names
* **[dave_resource_name_no_type_substring](#dave_resource_name_no_type_substring):** Avoid redundant information in resource names
* **[dave_variable_alphabetical_order](#dave_variable_alphabetical_order):** Sort variables alphabetically
* **[dave_variable_must_be_in_variables_file](#dave_variable_must_be_in_variables_file):** Only allow variables in variable.tf
* **[dave_variable_region][#dave_variable_region]:** Don't allow region as a variable

## Installation

### Building from Source

```bash
git clone https://github.com/yourusername/tflint-ruleset-dave-says.git
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
  version = "0.1.0"
  source  = "github.com/yourusername/tflint-ruleset-dave-says"
}
```

### Disabling Individual Rules

```hcl
rule "dave_label_min_length" {
  enabled = false
}

rule "dave_variable_alphabetical_order" {
  enabled = false
}
```

## Rule Details

### Label Rules

#### `dave_label_snake`
**Purpose**: Ensure labels use only use snake_case, or lowercase letters, numbers, and underscores.

**Examples**:
```hcl
# ✅ Valid

data "aws_s3_bucket" "existing_storage" {}
ephemeral "random_password" "db" {}
module "my_module" { source = "./modules/my_module" }
output "bucket_name" { value = aws_s3_bucket.user_data.bucket }
resource "aws_s3_bucket" "storage123" {}
resource "aws_s3_bucket" "user_data" {}
variable "storage_name" {}`,

# ❌ Invalid

resource "aws_s3_bucket" "myStorage" {}    # Contains uppercase
resource "aws_s3_bucket" "user_Data" {}    # Contains uppercase
resource "aws_s3_bucket" "user-data" {}    # Contains dash
```

#### `dave_label_min_length`
**Purpose**: Ensure labels are at least 3 characters long or a well known name.

`db` and `s3` are two commonly used well known names.

**Examples**:
```hcl
# ✅ Valid

data "aws_iam_policy_document" "s3" {}
resource "aws_s3_bucket" "abc" {}
resource "aws_s3_bucket" "user_data" {}

# ❌ Invalid

resource "aws_s3_bucket" "a" {}  # Too short
resource "aws_s3_bucket" "yz" {}   # Too short
```

#### `dave_resource_label_no_type_substring`
**Purpose**: Prevent labels from containing words that appear in the resource type.

**Examples**:
```hcl
# ✅ Valid

resource "aws_s3_bucket" "user_data" {}
resource "aws_iam_role" "application" {}

# ❌ Invalid

resource "aws_s3_bucket" "user_bucket" {}  # Contains "bucket"
resource "aws_iam_role" "admin_role" {}    # Contains "role"
resource "aws_s3_bucket" "s3_storage" {}   # Contains "s3"
```

### Resource Name Rules

#### `dave_resource_name_no_type_substring`
**Purpose**: Prevent resource `name` and `name_prefix` attributes from containing words that appear in the resource type.

**Examples**:
```hcl
# ✅ Valid

resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}

# ❌ Invalid

resource "aws_s3_bucket" "main" {
  name = "my-application-bucket"  # Contains "bucket"
}

resource "aws_iam_role" "main" {
  name = "admin-role-policy"  # Contains "role"
}
```

#### `dave_resource_name_kebab`
**Purpose**: Ensure resource `name` and `name_prefix` attributes use kebab-case, or only lowercase letters, numbers, and dashes.

**Examples**:
```hcl
# ✅ Valid

resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}

resource "aws_iam_role" "main" {
  name_prefix = "admin-role-"
}

# ❌ Invalid
resource "aws_s3_bucket" "main" {
  name = "my_application_storage"  # Contains underscores
}
```

### Variable Rules

#### `dave_variable_must_be_in_variables_file`
**Purpose**: Ensure all variable blocks are declared in `variables.tf` files only.

**Examples**:
```hcl
# ✅ Valid
# variables.tf
variable "storage_name" {
  type = string
}

# ❌ Invalid
# main.tf or any other file
variable "storage_name" {
  type = string
}
```

#### `dave_variable_alphabetical_order`
**Purpose**: Ensure variables are sorted alphabetically by name.

**Examples**:
```hcl
# ✅ Valid
 - alphabetical order
variable "alpha" {}
variable "beta" {}
variable "gamma" {}

# ❌ Invalid
 - not alphabetical
variable "beta" {}
variable "alpha" {}
variable "gamma" {}
```

#### `dave_variable_region`
**Purpose**: Prevents use of 'region' as a variable name.

**Examples**:
```hcl
# ✅ Valid
variable "alpha" {}
variable "beta" {}
variable "gamma" {}

# ❌ Invalid
variable "region" {}
```


### AWS Policy Rules

#### `dave_aws_policy_no_jsonencode`
**Purpose**: Ensure AWS IAM policies reference `aws_iam_policy_document` data sources instead of using `jsonencode()`.

**Examples**:
```hcl
# ✅ Valid

resource "aws_iam_role" "main" {
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

# ❌ Invalid

resource "aws_iam_role" "main" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [...]
  })
}

# ✅ Valid
 - Non-AWS resources can use jsonencode
resource "kubernetes_config_map" "main" {
  data = {
    config = jsonencode({ key = "value" })
  }
}
```

## Usage Examples

### Running TFLint with Dave Says Rules

```bash
# Run all rules
tflint

# Run with specific configuration
tflint --config .tflint.hcl

# Show only warnings from dave-says plugin
tflint --only=dave_label_snake --only=dave_variable_alphabetical_order
```

### Example Output

```
Warning: Variable 'storageVar' is not in alphabetical order (dave_variable_alphabetical_order)

  on test.tf line 23:
  23: variable "storageVar" {

Warning: Variable name 'storageVar' must contain only lowercase letters, numbers, and underscores (dave_label_snake)

  on test.tf line 23:
  23: variable "storageVar" {

Warning: Variable 'applicationVar' is not in alphabetical order. Expected position after 'applicationVar' (dave_variable_alphabetical_order)

  on test.tf line 27:
  27: variable "applicationVar" {

Warning: Variable name 'applicationVar' must contain only lowercase letters, numbers, and underscores (dave_label_snake)

  on test.tf line 27:
  27: variable "applicationVar" {
```

## Development

### Requirements

- Go 1.21 or newer
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