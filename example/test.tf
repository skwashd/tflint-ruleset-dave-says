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
