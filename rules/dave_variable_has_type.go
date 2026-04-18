package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveVariableHasTypeRule ensures all variables have an explicit type constraint.
type DaveVariableHasTypeRule struct {
	BaseRule
}

func NewDaveVariableHasTypeRule() *DaveVariableHasTypeRule {
	return &DaveVariableHasTypeRule{
		BaseRule: BaseRule{ruleName: "dave_variable_has_type"},
	}
}

func (r *DaveVariableHasTypeRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_variable_has_type.md"
}

func (r *DaveVariableHasTypeRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "type"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		_, exists := block.Body.Attributes["type"]
		if !exists {
			if err := EmitIssue(runner, r,
				fmt.Sprintf("Variable %q is missing a type constraint. Without one, Terraform cannot catch invalid inputs at plan time.", block.Labels[0]),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
