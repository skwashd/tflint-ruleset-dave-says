package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveNoVpcIdVariable_Valid(t *testing.T) {
	rule := NewDaveNoVpcIdVariableRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "subnet_ids" {
  type = list(string)
}

variable "my_vpc_id" {
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
}

func Test_DaveNoVpcIdVariable_Invalid(t *testing.T) {
	rule := NewDaveNoVpcIdVariableRule()

	runner := helper.TestRunner(t, map[string]string{
		"variables.tf": `
variable "vpc_id" {
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

	expected := `Do not use "vpc_id" as a variable. Derive it from the first subnet using data.aws_subnet to prevent mismatches between VPC ID and subnet IDs.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}
