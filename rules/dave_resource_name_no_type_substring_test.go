package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveResourceNameNoTypeSubstring(t *testing.T) {
	rule := NewDaveResourceNameNoTypeSubstringRule()

	// Test valid case - no type substring in name
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "my-application-storage"
}

resource "aws_iam_role" "admin" {
  name_prefix = "application-admin-"
}

resource "aws_security_group" "web" {
  name = "frontend-web-servers"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - name contains "bucket" substring
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "my-bucket-storage"
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
		expected := "Resource name attribute 'my-bucket-storage' contains substring 'bucket' from resource type 'aws_s3_bucket'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test invalid case - name_prefix contains "iam" substring
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_role" "admin" {
  name_prefix = "iam-admin-role-"
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
		expected := "Resource name_prefix attribute 'iam-admin-role-' contains substring 'iam' from resource type 'aws_iam_role'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test invalid case - name contains "security" substring
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "web" {
  name = "web-security-frontend"
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
		expected := "Resource name attribute 'web-security-frontend' contains substring 'security' from resource type 'aws_security_group'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test invalid case - name contains "group" substring
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "web" {
  name = "web-group-frontend"
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
		expected := "Resource name attribute 'web-group-frontend' contains substring 'group' from resource type 'aws_security_group'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test case - multiple violations (both name and name_prefix)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "s3-bucket-storage"
  name_prefix = "bucket-prefix-"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(runner.Issues))
	}

	// Test case - case insensitive matching
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  name = "my-Bucket-storage"
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
		expected := "Resource name attribute 'my-Bucket-storage' contains substring 'bucket' from resource type 'aws_s3_bucket'"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test case - skip evaluation errors (variables, references, etc.)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
variable "bucket_name" {
  type = string
}

resource "aws_s3_bucket" "main" {
  name = var.bucket_name
}

resource "aws_iam_role" "admin" {
  name_prefix = "${local.prefix}-role-"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for non-evaluable expressions, got %d", len(runner.Issues))
	}

	// Test case - no name or name_prefix attributes
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_s3_bucket" "main" {
  bucket = "my-bucket"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when name/name_prefix attributes are not present, got %d", len(runner.Issues))
	}
}