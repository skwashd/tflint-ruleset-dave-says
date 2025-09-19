package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// BaseRule provides common functionality for all Dave Says rules
type BaseRule struct {
	tflint.DefaultRule
	ruleName string
}

func (r *BaseRule) Name() string {
	return r.ruleName
}

func (r *BaseRule) Enabled() bool {
	return true
}

func (r *BaseRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *BaseRule) Link() string {
	return ""
}

// EmitIssue is a helper to emit issues with consistent formatting
func EmitIssue(runner tflint.Runner, rule tflint.Rule, message string, rng hcl.Range) error {
	return runner.EmitIssue(rule, message, rng)
}
