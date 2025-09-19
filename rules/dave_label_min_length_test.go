package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveLabelMinLength(t *testing.T) {
	rule := NewDaveLabelMinLengthRule()

	// Test valid case
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
data "aws_iam_policy_document" "s3" {}
resource "aws_secretsmanager_secret" "db" {}
resource "aws_s3_bucket" "abc" {}
resource "aws_s3_bucket" "user_data" {}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - only violates Rule 2 (too short)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `resource "aws_s3_bucket" "ab" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Resource label 'ab' must be at least 3 characters long or a well known name such as db or s3"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
