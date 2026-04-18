package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveIamNoInlinePolicy_Valid(t *testing.T) {
	rule := NewDaveIamNoInlinePolicyRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_policy" "app" {
  name   = "app-policy"
  policy = data.aws_iam_policy_document.app.json
}

resource "aws_iam_role_policy_attachment" "app" {
  role       = aws_iam_role.app.name
  policy_arn = aws_iam_policy.app.arn
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

func Test_DaveIamNoInlinePolicy_RolePolicy(t *testing.T) {
	rule := NewDaveIamNoInlinePolicyRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_role_policy" "app" {
  name   = "app-policy"
  role   = aws_iam_role.app.id
  policy = data.aws_iam_policy_document.app.json
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := "Use aws_iam_policy + aws_iam_role_policy_attachment instead of aws_iam_role_policy. Inline policies are not reusable and are harder to audit in the IAM console."
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveIamNoInlinePolicy_UserPolicy(t *testing.T) {
	rule := NewDaveIamNoInlinePolicyRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_user_policy" "app" {
  name   = "app-policy"
  user   = aws_iam_user.app.name
  policy = data.aws_iam_policy_document.app.json
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

func Test_DaveIamNoInlinePolicy_GroupPolicy(t *testing.T) {
	rule := NewDaveIamNoInlinePolicyRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_group_policy" "app" {
  name   = "app-policy"
  group  = aws_iam_group.app.name
  policy = data.aws_iam_policy_document.app.json
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
