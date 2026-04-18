package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveOutputMustBeInOutputsFile_Valid(t *testing.T) {
	rule := NewDaveOutputMustBeInOutputsFileRule()

	runner := helper.TestRunner(t, map[string]string{
		"outputs.tf": `
output "bucket_name" {
  value = aws_s3_bucket.main.bucket
}

output "bucket_arn" {
  value = aws_s3_bucket.main.arn
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

func Test_DaveOutputMustBeInOutputsFile_Invalid(t *testing.T) {
	rule := NewDaveOutputMustBeInOutputsFileRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
output "bucket_name" {
  value = aws_s3_bucket.main.bucket
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Output "bucket_name" is defined in main.tf. All outputs should be in outputs.tf for consistent file organization.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveOutputMustBeInOutputsFile_MixedFiles(t *testing.T) {
	rule := NewDaveOutputMustBeInOutputsFileRule()

	runner := helper.TestRunner(t, map[string]string{
		"outputs.tf": `
output "valid_output" {
  value = "ok"
}
`,
		"main.tf": `
output "misplaced_output" {
  value = "not ok"
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
