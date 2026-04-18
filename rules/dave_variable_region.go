package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveVariableRegionRule struct {
	BaseRule
}

func NewDaveVariableRegionRule() *DaveVariableRegionRule {
	return &DaveVariableRegionRule{
		BaseRule: BaseRule{ruleName: "dave_variable_region"},
	}
}

func (r *DaveVariableRegionRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_variable_region.md"
}

func (r *DaveVariableRegionRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		label := block.Labels[0]
		if label == "region" {
			err := EmitIssue(runner, r, fmt.Sprintf("Variable '%s' is not allowed; use provider default region", label), block.DefRange)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
