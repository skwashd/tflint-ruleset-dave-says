package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveCloudwatchLogRetentionRule ensures CloudWatch log groups have the correct retention_in_days.
// The expected value is configurable via the retention_days rule config (default: 30).
type DaveCloudwatchLogRetentionRule struct {
	BaseRule
	RetentionDays int
}

func NewDaveCloudwatchLogRetentionRule() *DaveCloudwatchLogRetentionRule {
	return &DaveCloudwatchLogRetentionRule{
		BaseRule:       BaseRule{ruleName: "dave_cloudwatch_log_retention"},
		RetentionDays: 30,
	}
}

func (r *DaveCloudwatchLogRetentionRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_cloudwatch_log_retention.md"
}

func (r *DaveCloudwatchLogRetentionRule) Check(runner tflint.Runner) error {
	cfg := struct {
		RetentionDays int `hclext:"retention_days,optional"`
	}{RetentionDays: r.RetentionDays}
	if err := runner.DecodeRuleConfig(r.Name(), &cfg); err != nil {
		return err
	}
	expected := cfg.RetentionDays

	resources, err := runner.GetResourceContent("aws_cloudwatch_log_group", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "retention_in_days"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attr, exists := resource.Body.Attributes["retention_in_days"]
		if !exists {
			if err := runner.EmitIssue(r,
				fmt.Sprintf("CloudWatch log group is missing retention_in_days. Set retention_in_days = %d.", expected),
				resource.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		var retention int
		if err := runner.EvaluateExpr(attr.Expr, &retention, nil); err != nil {
			continue
		}

		if retention != expected {
			if err := runner.EmitIssueWithFix(r,
				fmt.Sprintf("CloudWatch log group retention_in_days is %d, expected %d.", retention, expected),
				attr.Expr.Range(),
				func(f tflint.Fixer) error {
					return f.ReplaceText(attr.Expr.Range(), fmt.Sprintf("%d", expected))
				},
			); err != nil {
				return err
			}
		}
	}
	return nil
}
