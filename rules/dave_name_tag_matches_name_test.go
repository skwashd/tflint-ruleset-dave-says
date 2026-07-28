package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_DaveNameTagMatchesName_Match(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = "web"
  }
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

func Test_DaveNameTagMatchesName_MismatchWithFix(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = "wrong"
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

	expected := `Name tag value "wrong" does not match the resource's name attribute value "web"; they must be identical`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}

	helper.AssertChanges(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = "web"
  }
}
`,
	}, runner.Changes())
}

func Test_DaveNameTagMatchesName_MismatchWithResolvableVariable(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
variable "thing" {
  default = "web"
}

resource "aws_db_instance" "this" {
  name = var.thing
  tags = {
    Name = "wrong"
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

	// The fix mirrors the name attribute's source text, not its resolved
	// value, so a reference stays a reference after the fix.
	helper.AssertChanges(t, map[string]string{
		"main.tf": `
variable "thing" {
  default = "web"
}

resource "aws_db_instance" "this" {
  name = var.thing
  tags = {
    Name = var.thing
  }
}
`,
	}, runner.Changes())
}

func Test_DaveNameTagMatchesName_MergeWithLiteralName(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	t.Run("match", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
locals {
  common_tags = { Env = "prod" }
}

resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.common_tags, { Name = "web" })
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 0 {
			t.Errorf("expected no issues, got %d", len(runner.Issues))
		}
	})

	t.Run("mismatch is flagged and fixed", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
locals {
  common_tags = { Env = "prod" }
}

resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.common_tags, { Name = "wrong" })
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
  common_tags = { Env = "prod" }
}

resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.common_tags, { Name = "web" })
}
`,
		}, runner.Changes())
	})
}

func Test_DaveNameTagMatchesName_NestedMerge(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	t.Run("match with Name two levels deep", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.a, merge(local.b, { Name = "web" }))
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 0 {
			t.Errorf("expected no issues, got %d", len(runner.Issues))
		}
	})

	t.Run("mismatch two levels deep is flagged and fixed", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.a, merge(local.b, { Name = "wrong" }))
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
resource "aws_db_instance" "this" {
  name = "web"
  tags = merge(local.a, merge(local.b, { Name = "web" }))
}
`,
		}, runner.Changes())
	})

	t.Run("last write wins across nesting levels", func(t *testing.T) {
		// Terraform's merge() applies keys left to right, so later arguments
		// override earlier ones. Two Name entries exist here - one in the
		// outer merge's first argument, and one in the deepest nested merge's
		// last argument. The latter is the effective value and is the one
		// that should be compared and fixed.
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = merge({ Name = "first" }, merge({ Name = "second" }, { Name = "third" }))
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 1 {
			t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
		}

		expected := `Name tag value "third" does not match the resource's name attribute value "web"; they must be identical`
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}

		helper.AssertChanges(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = merge({ Name = "first" }, merge({ Name = "second" }, { Name = "web" }))
}
`,
		}, runner.Changes())
	})
}

func Test_DaveNameTagMatchesName_VariableTags(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	t.Run("match", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
variable "tags" {
  default = {
    Name = "web"
  }
}

resource "aws_db_instance" "this" {
  name = "web"
  tags = var.tags
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 0 {
			t.Errorf("expected no issues, got %d", len(runner.Issues))
		}
	})

	t.Run("mismatch is flagged but not fixed", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
variable "tags" {
  default = {
    Name = "wrong"
  }
}

resource "aws_db_instance" "this" {
  name = "web"
  tags = var.tags
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 1 {
			t.Fatalf("expected 1 issue, got %d", len(runner.Issues))
		}

		expected := `Name tag value "wrong" does not match the resource's name attribute value "web"; they must be identical`
		if runner.Issues[0].Message != expected {
			t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
		}

		if len(runner.Changes()) != 0 {
			t.Errorf("expected no autofix when the Name entry has no source location, got %d changed files", len(runner.Changes()))
		}
	})
}

func Test_DaveNameTagMatchesName_QuotedKey(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	t.Run("match", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    "Name" = "web"
  }
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 0 {
			t.Errorf("expected no issues, got %d", len(runner.Issues))
		}
	})

	t.Run("mismatch", func(t *testing.T) {
		runner := helper.TestRunner(t, map[string]string{
			"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    "Name" = "wrong"
  }
}
`,
		})

		if err := rule.Check(runner); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(runner.Issues) != 1 {
			t.Errorf("expected 1 issue, got %d", len(runner.Issues))
		}
	})
}

func Test_DaveNameTagMatchesName_BothReferencesIdenticalSource(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = var.a
  tags = {
    Name = var.a
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for identical reference source, got %d", len(runner.Issues))
	}
}

func Test_DaveNameTagMatchesName_BothReferencesDifferingSource(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = var.a
  tags = {
    Name = var.b
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

	expected := `Name tag (var.b) does not match the resource's name attribute (var.a); they should reference the same value`
	if runner.Issues[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, runner.Issues[0].Message)
	}

	if len(runner.Changes()) != 0 {
		t.Errorf("expected no autofix for a reference mismatch, got %d changed files", len(runner.Changes()))
	}
}

func Test_DaveNameTagMatchesName_MixedKnownAndReferenceSkipped(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Name = var.b
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when only one side is resolvable, got %d", len(runner.Issues))
	}
}

func Test_DaveNameTagMatchesName_NoNameTag(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
  tags = {
    Env = "prod"
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when there is no Name tag, got %d", len(runner.Issues))
	}
}

func Test_DaveNameTagMatchesName_NoNameAttribute(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_instance" "this" {
  tags = {
    Name = "web"
  }
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues for a resource without a name attribute, got %d", len(runner.Issues))
	}
}

func Test_DaveNameTagMatchesName_NoTags(t *testing.T) {
	rule := NewDaveNameTagMatchesNameRule()

	runner := helper.TestRunner(t, map[string]string{
		"main.tf": `
resource "aws_db_instance" "this" {
  name = "web"
}
`,
	})

	if err := rule.Check(runner); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if len(runner.Issues) != 0 {
		t.Errorf("expected no issues when there are no tags, got %d", len(runner.Issues))
	}
}
