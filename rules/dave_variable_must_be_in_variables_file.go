package rules

import (
	"fmt"
	"path/filepath"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveVariableMustBeInVariablesFileRule struct {
	BaseRule
}

func NewDaveVariableMustBeInVariablesFileRule() *DaveVariableMustBeInVariablesFileRule {
	return &DaveVariableMustBeInVariablesFileRule{
		BaseRule: BaseRule{ruleName: "dave_variable_must_be_in_variables_file"},
	}
}

func (r *DaveVariableMustBeInVariablesFileRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_variable_must_be_in_variables_file.md"
}

func (r *DaveVariableMustBeInVariablesFileRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) >= 1 {
			varName := block.Labels[0]
			filename := filepath.Base(block.DefRange.Filename)
			
			if filename != "variables.tf" {
				err := EmitIssue(runner, r, fmt.Sprintf("Variable '%s' must be declared in variables.tf, not in %s", varName, filename), block.DefRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
