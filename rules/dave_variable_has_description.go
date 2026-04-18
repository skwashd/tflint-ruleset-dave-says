package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveVariableHasDescriptionRule ensures all variables have a description.
type DaveVariableHasDescriptionRule struct {
	BaseRule
}

func NewDaveVariableHasDescriptionRule() *DaveVariableHasDescriptionRule {
	return &DaveVariableHasDescriptionRule{
		BaseRule: BaseRule{ruleName: "dave_variable_has_description"},
	}
}

func (r *DaveVariableHasDescriptionRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_variable_has_description.md"
}

func (r *DaveVariableHasDescriptionRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "variable",
				LabelNames: []string{"name"},
				Body: &hclext.BodySchema{
					Attributes: []hclext.AttributeSchema{
						{Name: "description"},
					},
				},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		_, exists := block.Body.Attributes["description"]
		if !exists {
			if err := EmitIssue(runner, r,
				fmt.Sprintf("Variable %q is missing a description. Without one, module consumers must read the implementation to understand what it controls.", block.Labels[0]),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
