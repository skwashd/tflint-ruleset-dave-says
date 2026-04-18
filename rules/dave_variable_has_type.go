package rules

import (
	"fmt"

	"github.com/hashicorp/hcl/v2/hclsyntax"
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
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}

	for _, file := range files {
		body, ok := file.Body.(*hclsyntax.Body)
		if !ok {
			continue
		}

		for _, block := range body.Blocks {
			if block.Type != "variable" || len(block.Labels) == 0 {
				continue
			}

			hasType := false
			for _, attr := range block.Body.Attributes {
				if attr.Name == "type" {
					hasType = true
					break
				}
			}

			if !hasType {
				if err := EmitIssue(runner, r,
					fmt.Sprintf("Variable %q is missing a type constraint. Without one, Terraform cannot catch invalid inputs at plan time.", block.Labels[0]),
					block.DefRange(),
				); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
