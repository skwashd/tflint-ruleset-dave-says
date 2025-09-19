package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveLabelNoTypeSubstring(t *testing.T) {
	rule := NewDaveLabelNoTypeSubstringRule()

	// Test valid case
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "user_data" {}
resource "aws_iam_role" "application" {}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - bucket substring
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `resource "aws_s3_bucket" "user_bucket" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Resource label 'user_bucket' contains substring 'bucket' from resource type 'aws_s3_bucket'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
