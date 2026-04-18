package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveVariableAlphabeticalOrder(t *testing.T) {
	rule := NewDaveVariableAlphabeticalOrderRule()

	// Test valid case - variables in alphabetical order
	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "application_name" {
  type = string
}

variable "environment" {
  type = string
}

variable "region" {
  type = string
}

variable "tags" {
  type = map(string)
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - variables out of alphabetical order
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "environment" {
  type = string
}

variable "application_name" {
  type = string
}

variable "region" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Variable 'environment' is not in alphabetical order"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test case - first variable out of order (no expected position)
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "zebra" {
  type = string
}

variable "apple" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Variable 'zebra' is not in alphabetical order"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}

	// Test case - multiple variables out of order
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "zebra" {
  type = string
}

variable "banana" {
  type = string
}

variable "apple" {
  type = string
}

variable "cherry" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 3 {
		t.Errorf("expected 3 issues, got %d", len(runner.Issues))
	}

	// Test case - single variable (should not trigger)
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "single_variable" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for single variable, got %d", len(runner.Issues))
	}

	// Test case - no variables (should not trigger)
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
# No variables in this file
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for empty file, got %d", len(runner.Issues))
	}

	// Test case - mixed files
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "environment" {
  type = string
}

variable "application_name" {
  type = string
}
`,
		"main.tf": `
variable "zebra" {
  type = string
}

variable "banana" {
  type = string
}

variable "apple" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 4 {
		t.Errorf("expected 4 issues, got %d", len(runner.Issues))
	}

	// Map iteration order is non-deterministic, so check that expected issues exist
	// rather than asserting on a specific index
	messages := make(map[string]bool)
	for _, issue := range runner.Issues {
		messages[issue.Message] = true
	}
	if !messages["Variable 'environment' is not in alphabetical order"] {
		t.Error("expected issue for 'environment' not in alphabetical order")
	}
	if !messages["Variable 'zebra' is not in alphabetical order"] {
		t.Error("expected issue for 'zebra' not in alphabetical order")
	}

	// Test case - variables with underscores and numbers (proper alphabetical sorting)
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "app_name" {
  type = string
}

variable "app_version" {
  type = string
}

variable "database_port" {
  type = number
}

variable "database_user" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for properly sorted variables with underscores, got %d", len(runner.Issues))
	}

	// Test case - case sensitivity (uppercase comes before lowercase in Go sorting)
	runner = helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "apple" {
  type = string
}

variable "Apple" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues for case sensitivity, got %d", len(runner.Issues))
	}

	if len(runner.Issues) > 0 {
		expected := "Variable 'apple' is not in alphabetical order"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
