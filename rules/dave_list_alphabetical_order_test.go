package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveListAlphabeticalOrder_NotConfigured(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["b", "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when unconfigured, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_SortedSingleLine(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["a", "b", "c"]
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

func Test_DaveListAlphabeticalOrder_UnsortedSingleLineWithFix(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["b", "a", "c"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	expected := "List assigned to 'tags' is not sorted alphabetically."
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}

	helper.AssertChanges(t, map[string]string{
		"main.tf": `
locals {
  tags = ["a", "b", "c"]
}
`,
	}, runner.Changes())
}

func Test_DaveListAlphabeticalOrder_VariableReferenceSkipped(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["b", var.thing, "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when a list has a non-static element, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_EmptyAndSingleElement(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"empty.tf": `
locals {
  tags = []
}
`,
		"single.tf": `
locals {
  tags = ["only"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for empty/single-element lists, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_AttributeNotInSetIgnored(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  other = ["b", "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for an unconfigured attribute, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_NonListAttributeIgnored(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = "hello"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for a non-list value, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_DuplicatesWithFix(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["b", "a", "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	helper.AssertChanges(t, map[string]string{
		"main.tf": `
locals {
  tags = ["a", "a", "b"]
}
`,
	}, runner.Changes())
}

func Test_DaveListAlphabeticalOrder_CaseSensitiveDefault(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	// Capitals sort before lowercase in byte order, so ["Z", "a"] is sorted.
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["Z", "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for byte-sorted list, got %d", len(runner.Issues))
	}

	// ["a", "Z"] is out of order in byte terms.
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["a", "Z"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue for byte-unsorted list, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_CaseInsensitive(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}
	rule.CaseInsensitive = true

	// Case-insensitively, ["a", "Z"] is sorted.
	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["a", "Z"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for case-insensitive sorted list, got %d", len(runner.Issues))
	}

	// ["Z", "a"] is out of order case-insensitively.
	runner = helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["Z", "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Errorf("expected 1 issue for case-insensitive unsorted list, got %d", len(runner.Issues))
	}
}

func Test_DaveListAlphabeticalOrder_MultilineNotFixed(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = [
    "b",
    "a",
  ]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Changes()) != 0 {
		t.Errorf("expected no autofix for a multiline list, got %d changed files", len(runner.Changes()))
	}
}

func Test_DaveListAlphabeticalOrder_InlineCommentNotFixed(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags = ["b", /* note */ "a"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
	}

	if len(runner.Changes()) != 0 {
		t.Errorf("expected no autofix for a list with a comment, got %d changed files", len(runner.Changes()))
	}
}

func Test_DaveListAlphabeticalOrder_NestedBlock(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"actions"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_iam_policy" "this" {
  statement {
    actions = ["s3:PutObject", "s3:GetObject"]
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 1 {
		t.Fatalf("expected 1 issue for the nested attribute, got %d", len(runner.Issues))
	}

	expected := "List assigned to 'actions' is not sorted alphabetically."
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}
}

func Test_DaveListAlphabeticalOrder_MultipleLists(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()
	rule.Attributes = []string{"tags", "ports"}

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags  = ["b", "a"]
  ports = ["https", "http"]
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Fatalf("expected 2 issues, got %d", len(runner.Issues))
	}

	// Map iteration order is non-deterministic, so assert on the set of
	// messages rather than on an index.
	messages := make(map[string]bool)
	for _, issue := range runner.Issues {
		messages[issue.Message] = true
	}
	if !messages["List assigned to 'tags' is not sorted alphabetically."] {
		t.Error("expected issue for unsorted 'tags' list")
	}
	if !messages["List assigned to 'ports' is not sorted alphabetically."] {
		t.Error("expected issue for unsorted 'ports' list")
	}
}

func Test_DaveListAlphabeticalOrder_ConfigViaTflintHcl(t *testing.T) {
	rule := NewDaveListAlphabeticalOrderRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
locals {
  tags  = ["b", "a"]
  ports = ["https", "http"]
  other = ["z", "a"]
}
`,
		".tflint.hcl": `
rule "dave_list_alphabetical_order" {
  enabled          = true
  attributes       = ["tags", "ports"]
  case_insensitive = true
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 2 {
		t.Fatalf("expected 2 issues (tags and ports), got %d", len(runner.Issues))
	}

	messages := make(map[string]bool)
	for _, issue := range runner.Issues {
		messages[issue.Message] = true
	}
	if !messages["List assigned to 'tags' is not sorted alphabetically."] {
		t.Error("expected issue for unsorted 'tags' list")
	}
	if !messages["List assigned to 'ports' is not sorted alphabetically."] {
		t.Error("expected issue for unsorted 'ports' list")
	}
	if messages["List assigned to 'other' is not sorted alphabetically."] {
		t.Error("unexpected issue for 'other' (not in configured attributes)")
	}
}
