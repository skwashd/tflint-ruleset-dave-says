package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveVariableHasType_Valid(t *testing.T) {
	rule := NewDaveVariableHasTypeRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
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
}

func Test_DaveVariableHasType_Missing(t *testing.T) {
	rule := NewDaveVariableHasTypeRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
  description = "The name"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Variable "name" is missing a type constraint. Without one, Terraform cannot catch invalid inputs at plan time.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveVariableHasType_MixedPresence(t *testing.T) {
	rule := NewDaveVariableHasTypeRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "name" {
  type = string
}

variable "tags" {
  description = "Resource tags"
}

variable "enabled" {}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Errorf("expected 2 issues, got %d", len(runner.Issues))
	}
}
