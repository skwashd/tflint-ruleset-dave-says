package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveSecurityGroupNoInlineRulesRule flags inline ingress/egress blocks on aws_security_group.
// Use aws_vpc_security_group_ingress_rule and aws_vpc_security_group_egress_rule instead.
type DaveSecurityGroupNoInlineRulesRule struct {
	BaseRule
}

func NewDaveSecurityGroupNoInlineRulesRule() *DaveSecurityGroupNoInlineRulesRule {
	return &DaveSecurityGroupNoInlineRulesRule{
		BaseRule: BaseRule{ruleName: "dave_security_group_no_inline_rules"},
	}
}

func (r *DaveSecurityGroupNoInlineRulesRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_security_group_no_inline_rules.md"
}

var inlineRuleBlocks = map[string]string{
	"ingress": "aws_vpc_security_group_ingress_rule",
	"egress":  "aws_vpc_security_group_egress_rule",
}

func (r *DaveSecurityGroupNoInlineRulesRule) Check(runner tflint.Runner) error {
	blockSchemas := make([]hclext.BlockSchema, 0, len(inlineRuleBlocks))
	for blockType := range inlineRuleBlocks {
		blockSchemas = append(blockSchemas, hclext.BlockSchema{Type: blockType})
	}

	resources, err := runner.GetResourceContent("aws_security_group", &hclext.BodySchema{
		Blocks: blockSchemas,
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			replacement, ok := inlineRuleBlocks[block.Type]
			if !ok {
				continue
			}
			if err := runner.EmitIssue(r,
				fmt.Sprintf("Inline %q block on aws_security_group causes rule conflicts and is harder to manage. Use %s resources instead.", block.Type, replacement),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
