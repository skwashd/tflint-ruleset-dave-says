package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveS3NoPublicAcl_Valid(t *testing.T) {
	rule := NewDaveS3NoPublicAclRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket_acl" "data" {
  bucket = aws_s3_bucket.data.id
  acl    = "private"
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

func Test_DaveS3NoPublicAcl_PublicRead(t *testing.T) {
	rule := NewDaveS3NoPublicAclRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket_acl" "data" {
  bucket = aws_s3_bucket.data.id
  acl    = "public-read"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `S3 bucket ACL is "public-read". Public buckets are a leading cause of data breaches. Use CloudFront for web access.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveS3NoPublicAcl_PublicReadWrite(t *testing.T) {
	rule := NewDaveS3NoPublicAclRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket_acl" "data" {
  bucket = aws_s3_bucket.data.id
  acl    = "public-read-write"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}
}

func Test_DaveS3NoPublicAcl_AuthenticatedRead(t *testing.T) {
	rule := NewDaveS3NoPublicAclRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket_acl" "data" {
  bucket = aws_s3_bucket.data.id
  acl    = "authenticated-read"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}
}

func Test_DaveS3NoPublicAcl_NoAclAttribute(t *testing.T) {
	rule := NewDaveS3NoPublicAclRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket_acl" "data" {
  bucket = aws_s3_bucket.data.id
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
