package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveVariableMustBeInVariablesFile(t *testing.T) {
	rule := NewDaveVariableMustBeInVariablesFileRule()

	// Test valid case - variable in variables.tf
	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "storage_name" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues, got %d", len(runner.Issues))
	}

	// Test invalid case - variable in main.tf (only violates Rule 6)
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
variable "storage_name" {
  type = string
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
		expected := "Variable 'storage_name' must be declared in variables.tf, not in main.tf"
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}
	}
}
