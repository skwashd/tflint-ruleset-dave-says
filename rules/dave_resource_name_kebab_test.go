package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveResourceNameKebab(t *testing.T) {
	rule := NewDaveResourceNameKebabRule()

	// Test valid case - dashes only
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}

resource "aws_iam_role" "main" {
  name_prefix = "admin-role-"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - contains underscore (only violates Rule 10)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "my_application_storage"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Resource name attribute 'my_application_storage' must contain only lowercase letters, numbers, and dashes"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
