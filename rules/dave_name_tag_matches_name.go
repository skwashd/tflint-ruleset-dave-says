package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveNameTagMatchesNameRule flags resources whose `Name` tag does not match
// their `name` attribute. The `name` attribute is treated as the source of
// truth. The rule does not apply when a resource has no `name` attribute
// (e.g. aws_instance, which has no name argument at all), or no `Name` tag.
type DaveNameTagMatchesNameRule struct {
	BaseRule
}

func NewDaveNameTagMatchesNameRule() *DaveNameTagMatchesNameRule {
	return &DaveNameTagMatchesNameRule{
		BaseRule: BaseRule{ruleName: "dave_name_tag_matches_name"},
	}
}

func (r *DaveNameTagMatchesNameRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_name_tag_matches_name.md"
}

func (r *DaveNameTagMatchesNameRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "name"},
					{Name: "tags"},
				},
			}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		nameAttr, hasName := block.Body.Attributes["name"]
		if !hasName {
			continue // no name attribute -> rule does not apply
		}
		tagsAttr, hasTags := block.Body.Attributes["tags"]
		if !hasTags {
			continue // no tags -> nothing to compare
		}

		if err := r.checkResource(runner, nameAttr, tagsAttr); err != nil {
			return err
		}
	}

	return nil
}

// checkResource compares a single resource's name attribute against its Name
// tag, in two tiers: a precise comparison when the Name tag is a literal
// entry (which also enables autofix), and a value-only fallback when tags is
// some other expression (var.tags, a fully-static merge(), JSON) that still
// evaluates as a whole.
func (r *DaveNameTagMatchesNameRule) checkResource(runner tflint.Runner, nameAttr, tagsAttr *hclext.Attribute) error {
	if tagExpr, ok := findNameTagValueExpr(tagsAttr.Expr); ok {
		return r.comparePrecise(runner, nameAttr, tagExpr)
	}

	tagValue, ok := r.evalNameFromMap(runner, tagsAttr.Expr)
	if !ok {
		return nil // no Name entry could be located or evaluated -> skip
	}

	nameValue, ok := r.evalString(runner, nameAttr.Expr)
	if !ok {
		// name is an unresolvable reference; comparing it against a resolved
		// tag value with no source location is too ambiguous to flag safely.
		return nil
	}

	if nameValue == tagValue {
		return nil
	}

	return runner.EmitIssue(r, valueMismatchMessage(tagValue, nameValue), tagsAttr.Range)
}

// comparePrecise handles the case where the Name tag's value expression was
// located exactly, so both a precise issue location and an autofix (when the
// name value is known) are available.
func (r *DaveNameTagMatchesNameRule) comparePrecise(runner tflint.Runner, nameAttr *hclext.Attribute, tagExpr hcl.Expression) error {
	nameValue, nameEval := r.evalString(runner, nameAttr.Expr)
	tagValue, tagEval := r.evalString(runner, tagExpr)

	switch {
	case nameEval && tagEval:
		if nameValue == tagValue {
			return nil
		}
		message := valueMismatchMessage(tagValue, nameValue)
		if nameSrc, ok := r.exprSource(runner, nameAttr.Expr); ok {
			return runner.EmitIssueWithFix(r, message, tagExpr.Range(),
				func(f tflint.Fixer) error {
					return f.ReplaceText(tagExpr.Range(), nameSrc)
				},
			)
		}
		return runner.EmitIssue(r, message, tagExpr.Range())

	case !nameEval && !tagEval:
		// Neither side is statically known, so both are references (e.g.
		// var.x). Compare source text: identical source is provably
		// consistent. Differing source might still resolve to the same value
		// at apply time, so this can produce an occasional false positive —
		// but silently skipping every reference-based Name tag would miss
		// real drift like `name = var.a` vs `Name = var.b`.
		nameSrc, nameSrcOk := r.exprSource(runner, nameAttr.Expr)
		tagSrc, tagSrcOk := r.exprSource(runner, tagExpr)
		if !nameSrcOk || !tagSrcOk {
			return nil
		}
		nameSrc, tagSrc = strings.TrimSpace(nameSrc), strings.TrimSpace(tagSrc)
		if nameSrc == tagSrc {
			return nil
		}
		message := fmt.Sprintf(
			"Name tag (%s) does not match the resource's name attribute (%s); they should reference the same value",
			tagSrc, nameSrc,
		)
		return runner.EmitIssue(r, message, tagExpr.Range())

	default:
		// One side resolved to a known value and the other is an
		// unresolvable reference; too ambiguous to compare safely.
		return nil
	}
}

func valueMismatchMessage(tagValue, nameValue string) string {
	return fmt.Sprintf(
		"Name tag value %q does not match the resource's name attribute value %q; they must be identical",
		tagValue, nameValue,
	)
}

// findNameTagValueExpr recursively locates the value expression of a literal
// `Name` entry: directly in an object literal, or inside merge()/parentheses.
// merge() arguments are scanned in reverse, since later arguments win when
// Terraform merges maps. Returns false when no literal Name entry exists (for
// example tags = var.tags, or a Name key that is itself dynamic).
func findNameTagValueExpr(expr hcl.Expression) (hcl.Expression, bool) {
	switch e := expr.(type) {
	case *hclsyntax.ObjectConsExpr:
		for _, item := range e.Items {
			// ObjectConsKeyExpr.Value(nil) resolves a bare key (`Name = ...`)
			// or a quoted key (`"Name" = ...`) to cty.StringVal("Name") with
			// no evaluation context; a dynamic key returns diagnostics and is
			// skipped.
			keyVal, diags := item.KeyExpr.Value(nil)
			if diags.HasErrors() || keyVal.IsNull() || !keyVal.IsKnown() || keyVal.Type() != cty.String {
				continue
			}
			if keyVal.AsString() == "Name" {
				return item.ValueExpr, true
			}
		}
	case *hclsyntax.FunctionCallExpr:
		if e.Name == "merge" {
			for i := len(e.Args) - 1; i >= 0; i-- {
				if v, ok := findNameTagValueExpr(e.Args[i]); ok {
					return v, true
				}
			}
		}
	case *hclsyntax.ParenthesesExpr:
		return findNameTagValueExpr(e.Expression)
	}
	return nil, false
}

// evalString evaluates expr to a known string. It returns false when the
// expression cannot be evaluated (references, unknown vars) or when
// evaluation panics on a nil-typed expression, mirroring the recover guard in
// dave_resource_name_no_type_substring.go.
func (r *DaveNameTagMatchesNameRule) evalString(runner tflint.Runner, expr hcl.Expression) (val string, ok bool) {
	defer func() {
		if recover() != nil {
			val, ok = "", false
		}
	}()

	var s string
	if err := runner.EvaluateExpr(expr, &s, nil); err != nil {
		return "", false
	}
	return s, true
}

// evalNameFromMap evaluates the whole tags expression to a map[string]string
// and returns its Name entry. This is the fallback used when tags isn't a
// literal object we can locate a Name entry in directly (var.tags, a
// fully-static merge(), JSON) but might still evaluate as a whole.
func (r *DaveNameTagMatchesNameRule) evalNameFromMap(runner tflint.Runner, expr hcl.Expression) (val string, ok bool) {
	defer func() {
		if recover() != nil {
			val, ok = "", false
		}
	}()

	m := map[string]string{}
	if err := runner.EvaluateExpr(expr, &m, nil); err != nil {
		return "", false
	}
	v, present := m["Name"]
	return v, present
}

// exprSource returns the raw source text spanning expr's range, used to
// mirror the name attribute into an autofix and to compare unresolvable
// reference expressions textually.
func (r *DaveNameTagMatchesNameRule) exprSource(runner tflint.Runner, expr hcl.Expression) (string, bool) {
	rng := expr.Range()

	files, err := runner.GetFiles()
	if err != nil {
		return "", false
	}

	file, ok := files[rng.Filename]
	if !ok || file == nil || rng.Start.Byte < 0 || rng.End.Byte > len(file.Bytes) || rng.Start.Byte > rng.End.Byte {
		return "", false
	}

	return string(file.Bytes[rng.Start.Byte:rng.End.Byte]), true
}
