package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveS3NoInlineConfig_Valid(t *testing.T) {
	rule := NewDaveS3NoInlineConfigRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}

resource "aws_s3_bucket_versioning" "data" {
  bucket = aws_s3_bucket.data.id
  versioning_configuration {
    status = "Enabled"
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}
}

func Test_DaveS3NoInlineConfig_InlineBlock(t *testing.T) {
	rule := NewDaveS3NoInlineConfigRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"

  versioning {
    enabled = true
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Inline "versioning" block is deprecated since AWS provider v4. Use the aws_s3_bucket_versioning resource instead.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveS3NoInlineConfig_InlineAttribute(t *testing.T) {
	rule := NewDaveS3NoInlineConfigRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
  acl    = "private"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Inline "acl" attribute is deprecated since AWS provider v4. Use the aws_s3_bucket_acl resource instead.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveS3NoInlineConfig_MultipleViolations(t *testing.T) {
	rule := NewDaveS3NoInlineConfigRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
  acl    = "private"

  versioning {
    enabled = true
  }

  logging {
    target_bucket = "logs-bucket"
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 3 {
		t.Errorf("expected 3 issues, got %d", len(runner.Issues))
	}
}
