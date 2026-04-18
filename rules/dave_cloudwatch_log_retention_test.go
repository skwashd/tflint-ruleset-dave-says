package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveCloudwatchLogRetention_Valid(t *testing.T) {
	rule := NewDaveCloudwatchLogRetentionRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_cloudwatch_log_group" "app" {
  name              = "/app/logs"
  retention_in_days = 30
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

func Test_DaveCloudwatchLogRetention_MissingRetention(t *testing.T) {
	rule := NewDaveCloudwatchLogRetentionRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_cloudwatch_log_group" "app" {
  name = "/app/logs"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := "CloudWatch log group is missing retention_in_days. Set retention_in_days = 30."
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveCloudwatchLogRetention_WrongValue(t *testing.T) {
	rule := NewDaveCloudwatchLogRetentionRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_cloudwatch_log_group" "app" {
  name              = "/app/logs"
  retention_in_days = 90
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := "CloudWatch log group retention_in_days is 90, expected 30."
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveCloudwatchLogRetention_NonCloudwatchIgnored(t *testing.T) {
	rule := NewDaveCloudwatchLogRetentionRule()

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

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}
}
