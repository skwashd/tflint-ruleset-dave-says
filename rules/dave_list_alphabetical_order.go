package rules

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/zclconf/go-cty/cty"

	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveListAlphabeticalOrderRule flags list literals whose string elements are
// not sorted alphabetically, for a configurable set of attribute names.
//
// The rule is opt-in: there is no universally-correct set of lists to sort, so
// it is a no-op unless the attributes config is set. By naming an attribute the
// user is asserting that element order is not semantically significant for it.
type DaveListAlphabeticalOrderRule struct {
	BaseRule
	Attributes      []string
	CaseInsensitive bool
}

func NewDaveListAlphabeticalOrderRule() *DaveListAlphabeticalOrderRule {
	return &DaveListAlphabeticalOrderRule{
		BaseRule: BaseRule{ruleName: "dave_list_alphabetical_order"},
	}
}

func (r *DaveListAlphabeticalOrderRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_list_alphabetical_order.md"
}

func (r *DaveListAlphabeticalOrderRule) Check(runner tflint.Runner) error {
	cfg := struct {
		Attributes      []string `hclext:"attributes,optional"`
		CaseInsensitive bool     `hclext:"case_insensitive,optional"`
	}{}
	if err := runner.DecodeRuleConfig(r.Name(), &cfg); err != nil {
		return err
	}
	if len(cfg.Attributes) > 0 {
		r.Attributes = cfg.Attributes
		r.CaseInsensitive = cfg.CaseInsensitive
	}

	// Opt-in by design: with nothing configured the rule does nothing.
	if len(r.Attributes) == 0 {
		return nil
	}

	wanted := make(map[string]struct{}, len(r.Attributes))
	for _, name := range r.Attributes {
		wanted[name] = struct{}{}
	}

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	// Comment ranges are lexed lazily, once per file, and cached.
	commentCache := make(map[string][]hcl.Range)

	for filename, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			// Non-native (e.g. JSON) body; the token-level checks below do
			// not apply, so skip it.
			continue
		}

		if err := r.checkBody(runner, body, wanted, filename, file, commentCache); err != nil {
			return err
		}
	}

	return nil
}

// checkBody recurses through a native body, checking matched attributes and
// descending into every nested block so attributes nested inside blocks are
// caught too.
func (r *DaveListAlphabeticalOrderRule) checkBody(runner tflint.Runner, body *hclsyntax.Body, wanted map[string]struct{}, filename string, file *hcl.File, commentCache map[string][]hcl.Range) error {
	for name, attr := range body.Attributes {
		if _, ok := wanted[name]; !ok {
			continue
		}
		if err := r.checkAttribute(runner, attr, filename, file, commentCache); err != nil {
			return err
		}
	}

	for _, block := range body.Blocks {
		if err := r.checkBody(runner, block.Body, wanted, filename, file, commentCache); err != nil {
			return err
		}
	}

	return nil
}

func (r *DaveListAlphabeticalOrderRule) checkAttribute(runner tflint.Runner, attr *hclsyntax.Attribute, filename string, file *hcl.File, commentCache map[string][]hcl.Range) error {
	// Only act on `[...]` source (list/set/tuple syntax). Anything else — a
	// string, a function call, a variable ref — is ignored.
	tuple, ok := attr.Expr.(*hclsyntax.TupleConsExpr)
	if !ok {
		return nil
	}

	// Lists shorter than two elements are trivially sorted.
	if len(tuple.Exprs) < 2 {
		return nil
	}

	// Extract a comparison key per element. An element is static only if it
	// evaluates to a known string with no errors. If any element is non-static
	// (variable ref, interpolation, etc.) the ordering is ambiguous, so skip
	// the whole list.
	keys := make([]string, len(tuple.Exprs))
	for i, elem := range tuple.Exprs {
		v, diags := elem.Value(nil)
		if diags.HasErrors() || !v.IsKnown() || v.Type() != cty.String {
			return nil
		}
		keys[i] = r.comparisonKey(v.AsString())
	}

	if r.isSorted(keys) {
		return nil
	}

	message := fmt.Sprintf("List assigned to '%s' is not sorted alphabetically.", attr.Name)

	if r.canFix(tuple, filename, file, commentCache) {
		fixed := r.buildFixed(tuple, keys, file)
		return runner.EmitIssueWithFix(r, message, tuple.Range(),
			func(f tflint.Fixer) error {
				return f.ReplaceText(tuple.Range(), fixed)
			},
		)
	}

	return runner.EmitIssue(r, message, tuple.Range())
}

// comparisonKey returns the value used for ordering, lower-cased when
// case_insensitive is set.
func (r *DaveListAlphabeticalOrderRule) comparisonKey(s string) string {
	if r.CaseInsensitive {
		return strings.ToLower(s)
	}
	return s
}

// isSorted reports whether keys already match a sorted copy of themselves.
func (r *DaveListAlphabeticalOrderRule) isSorted(keys []string) bool {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)

	for i := range keys {
		if keys[i] != sorted[i] {
			return false
		}
	}
	return true
}

// canFix reports whether the list can be safely rewritten: it must be on a
// single line and contain no comments.
func (r *DaveListAlphabeticalOrderRule) canFix(tuple *hclsyntax.TupleConsExpr, filename string, file *hcl.File, commentCache map[string][]hcl.Range) bool {
	rng := tuple.Range()
	if rng.Start.Line != rng.End.Line {
		return false
	}

	for _, c := range r.fileComments(filename, file, commentCache) {
		if c.Start.Byte >= rng.Start.Byte && c.End.Byte <= rng.End.Byte {
			return false
		}
	}

	return true
}

// fileComments lexes the file once and returns the ranges of all comment
// tokens, caching the result per file. Tokens are used rather than substring
// matching so a '#' inside a string value is not mistaken for a comment.
func (r *DaveListAlphabeticalOrderRule) fileComments(filename string, file *hcl.File, cache map[string][]hcl.Range) []hcl.Range {
	if ranges, ok := cache[filename]; ok {
		return ranges
	}

	var ranges []hcl.Range
	tokens, _ := hclsyntax.LexConfig(file.Bytes, filename, hcl.InitialPos)
	for _, tok := range tokens {
		if tok.Type == hclsyntax.TokenComment {
			ranges = append(ranges, tok.Range)
		}
	}

	cache[filename] = ranges
	return ranges
}

// buildFixed produces the sorted single-line replacement. Each element keeps
// its original source slice so exact quoting and escaping are preserved.
func (r *DaveListAlphabeticalOrderRule) buildFixed(tuple *hclsyntax.TupleConsExpr, keys []string, file *hcl.File) string {
	type element struct {
		key string
		src string
	}

	elements := make([]element, len(tuple.Exprs))
	for i, elem := range tuple.Exprs {
		rng := elem.Range()
		elements[i] = element{
			key: keys[i],
			src: string(file.Bytes[rng.Start.Byte:rng.End.Byte]),
		}
	}

	sort.SliceStable(elements, func(i, j int) bool {
		return elements[i].key < elements[j].key
	})

	parts := make([]string, len(elements))
	for i, e := range elements {
		parts[i] = e.src
	}

	return "[" + strings.Join(parts, ", ") + "]"
}
