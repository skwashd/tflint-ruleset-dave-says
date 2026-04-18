package rules

import (
	"fmt"
	"math/big"

	"github.com/zclconf/go-cty/cty"

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

// ApplyRuleConfig decodes per-rule configuration from .tflint.hcl.
//
//	rule "dave_cloudwatch_log_retention" {
//	  enabled        = true
//	  retention_days = 14
//	}
func (r *DaveCloudwatchLogRetentionRule) ApplyRuleConfig(body *hclext.BodyContent) error {
	attr, exists := body.Attributes["retention_days"]
	if !exists {
		return nil
	}

	val, diags := attr.Expr.Value(nil)
	if diags.HasErrors() {
		return fmt.Errorf("evaluating retention_days: %s", diags.Error())
	}

	if val.Type() != cty.Number {
		return fmt.Errorf("retention_days must be a number, got %s", val.Type().FriendlyName())
	}

	bf := val.AsBigFloat()
	i, accuracy := bf.Int64()
	if accuracy != big.Exact {
		return fmt.Errorf("retention_days must be an exact integer")
	}

	r.RetentionDays = int(i)
	return nil
}

func (r *DaveCloudwatchLogRetentionRule) Check(runner tflint.Runner) error {
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
				fmt.Sprintf("CloudWatch log group is missing retention_in_days. Set retention_in_days = %d.", r.RetentionDays),
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

		if retention != r.RetentionDays {
			expected := r.RetentionDays
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
