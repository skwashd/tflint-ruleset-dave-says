package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveNoVpcIdVariableRule flags variables named "vpc_id". VPC ID should be
// derived from subnet data sources to prevent mismatches.
type DaveNoVpcIdVariableRule struct {
	BaseRule
}

func NewDaveNoVpcIdVariableRule() *DaveNoVpcIdVariableRule {
	return &DaveNoVpcIdVariableRule{
		BaseRule: BaseRule{ruleName: "dave_no_vpc_id_variable"},
	}
}

func (r *DaveNoVpcIdVariableRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_no_vpc_id_variable.md"
}

func (r *DaveNoVpcIdVariableRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if block.Labels[0] == "vpc_id" {
			if err := EmitIssue(runner, r,
				`Do not use "vpc_id" as a variable. Derive it from the first subnet using data.aws_subnet to prevent mismatches between VPC ID and subnet IDs.`,
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
