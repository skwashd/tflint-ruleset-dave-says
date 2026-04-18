package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveS3BucketNamespace_Valid(t *testing.T) {
	rule := NewDaveS3BucketNamespaceRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket_namespace = "account-regional"
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

func Test_DaveS3BucketNamespace_Missing(t *testing.T) {
	rule := NewDaveS3BucketNamespaceRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket = "my-bucket"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `S3 buckets should use bucket_namespace = "account-regional" to avoid global name collisions and confused-deputy risks. Requires AWS provider >= 6.37.0.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveS3BucketNamespace_WrongValue(t *testing.T) {
	rule := NewDaveS3BucketNamespaceRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "data" {
  bucket_namespace = "global"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `S3 bucket_namespace should be "account-regional", not "global"`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}
