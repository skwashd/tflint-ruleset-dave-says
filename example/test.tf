# Test file with isolated rule violations

resource "aws_s3_bucket" "myStorage" {
  name = "my-application-storage"
}

resource "aws_s3_bucket" "ab" {}

resource "aws_s3_bucket" "user-data" {}

resource "aws_s3_bucket" "user_bucket" {}

resource "aws_iam_role" "main" {
  name = "admin-role-policy"
}

resource "aws_iam_role" "application" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
  })
}

variable "storageVar" {
  type = string
}

variable "applicationVar" {
  type = string
}

# New rule violations below

# dave_iam_no_inline_policy
resource "aws_iam_role_policy" "inline" {
  name   = "inline-policy"
  role   = aws_iam_role.main.id
  policy = data.aws_iam_policy_document.app.json
}

# dave_no_vpc_id_variable
variable "vpc_id" {
  type        = string
  description = "VPC ID"
}

# dave_cloudwatch_log_retention (missing retention_in_days)
resource "aws_cloudwatch_log_group" "app" {
  name = "/app/logs"
}

# dave_s3_bucket_namespace (missing bucket_namespace)
resource "aws_s3_bucket" "namespaced" {
  bucket = "my-data-bucket"
}

# dave_s3_no_inline_config (deprecated inline versioning block)
resource "aws_s3_bucket" "inline_config" {
  bucket = "inline-config-bucket"

  versioning {
    enabled = true
  }
}

# dave_s3_no_public_acl
resource "aws_s3_bucket_acl" "public" {
  bucket = aws_s3_bucket.namespaced.id
  acl    = "public-read"
}

# dave_security_group_no_inline_rules
resource "aws_security_group" "inline" {
  name = "inline-sg"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# dave_output_must_be_in_outputs_file (output in wrong file)
output "misplaced_output" {
  value = aws_s3_bucket.namespaced.arn
}

# dave_variable_has_description (missing description)
variable "no_description" {
  type = string
}

# dave_variable_has_type (missing type)
variable "no_type" {
  description = "A variable without a type"
}
