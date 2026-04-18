package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveSecurityGroupNoInlineRules_Valid(t *testing.T) {
	rule := NewDaveSecurityGroupNoInlineRulesRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "app" {
  name        = "app-sg"
  description = "Application security group"
  vpc_id      = aws_vpc.main.id
}

resource "aws_vpc_security_group_ingress_rule" "app" {
  security_group_id = aws_security_group.app.id
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
  cidr_ipv4         = "0.0.0.0/0"
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

func Test_DaveSecurityGroupNoInlineRules_InlineIngress(t *testing.T) {
	rule := NewDaveSecurityGroupNoInlineRulesRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "app" {
  name = "app-sg"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Inline "ingress" block on aws_security_group causes rule conflicts and is harder to manage. Use aws_vpc_security_group_ingress_rule resources instead.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveSecurityGroupNoInlineRules_InlineEgress(t *testing.T) {
	rule := NewDaveSecurityGroupNoInlineRulesRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "app" {
  name = "app-sg"

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := `Inline "egress" block on aws_security_group causes rule conflicts and is harder to manage. Use aws_vpc_security_group_egress_rule resources instead.`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveSecurityGroupNoInlineRules_BothIngressAndEgress(t *testing.T) {
	rule := NewDaveSecurityGroupNoInlineRulesRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_security_group" "app" {
  name = "app-sg"

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
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
