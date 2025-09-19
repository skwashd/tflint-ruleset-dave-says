package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveLabelLowercaseOnly(t *testing.T) {
	rule := NewDaveLabelSnakeRule()

	// Test valid case
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
data "aws_s3_bucket" "existing_storage" {}
ephemeral "random_password" "db" {}
module "my_module" { source = "./modules/my_module" }
output "bucket_name" { value = aws_s3_bucket.user_data.bucket }
resource "aws_s3_bucket" "storage123" {}
resource "aws_s3_bucket" "user_data" {}
variable "storage_name" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Failing resource
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `resource "aws_s3_bucket" "myStorage" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Resource label 'myStorage' must contain only lowercase letters, numbers, and underscores"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Failing Variable
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `variable "storageVar" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Variable name 'storageVar' must contain only lowercase letters, numbers, and underscores"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Failing data source with dash
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `data "aws_iam_policy_document" "assume-role" {}`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Data label 'assume-role' must contain only lowercase letters, numbers, and underscores"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
