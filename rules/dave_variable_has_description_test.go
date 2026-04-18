package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveVariableHasDescription_Valid(t *testing.T) {
	rule := NewDaveVariableHasDescriptionRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
  type        = string
  description = "The name of the resource"
}

variable "tags" {
  type        = map(string)
  description = "Tags to apply"
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

func Test_DaveVariableHasDescription_Missing(t *testing.T) {
	rule := NewDaveVariableHasDescriptionRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
  type = string
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Variable "name" is missing a description. Without one, module consumers must read the implementation to understand what it controls.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveVariableHasDescription_MixedPresence(t *testing.T) {
	rule := NewDaveVariableHasDescriptionRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
  type        = string
  description = "The name"
}

variable "tags" {
  type = map(string)
}

variable "enabled" {
  type = bool
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(runner.Issues))
	}
}
