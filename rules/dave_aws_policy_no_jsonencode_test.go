package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveAwsPolicyNoJsonencode(t *testing.T) {
	rule := NewDaveAwsPolicyNoJsonencodeRule()

	// Test valid case - using data source
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_role" "main" {
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - using jsonencode
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_role" "main" {
  assume_role_policy = jsonencode({
    Version = "2012-10-17"
  })
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
		expected := "AWS resource 'aws_iam_role' attribute 'assume_role_policy' should reference an aws_iam_policy_document data source instead of using jsonencode()"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test non-AWS resource (should be ignored)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "kubernetes_config_map" "main" {
  data = {
    config = jsonencode({
      key = "value"
    })
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for non-AWS resource, got %d", len(runner.Issues))
	}
}
