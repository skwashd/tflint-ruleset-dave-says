package rules

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveVariableAlphabeticalOrderRule struct {
	BaseRule
}

func NewDaveVariableAlphabeticalOrderRule() *DaveVariableAlphabeticalOrderRule {
	return &DaveVariableAlphabeticalOrderRule{
		BaseRule: BaseRule{ruleName: "dave_variable_alphabetical_order"},
	}
}

func (r *DaveVariableAlphabeticalOrderRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_variable_alphabetical_order.md"
}

func (r *DaveVariableAlphabeticalOrderRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}, nil)
	if err != nil {
		return err
	}

	// Group variables by file
	fileVars := make(map[string][]*hclext.Block)
	for _, block := range content.Blocks {
		filename := filepath.Base(block.DefRange.Filename)
		fileVars[filename] = append(fileVars[filename], block)
	}

	// Check alphabetical order for each variables.tf file
	for _, blocks := range fileVars {
		if len(blocks) <= 1 {
			continue
		}

		// Extract variable names in order
		var names []string
		for _, block := range blocks {
			if len(block.Labels) >= 1 {
				names = append(names, block.Labels[0])
			}
		}

		// Create sorted copy
		sortedNames := make([]string, len(names))
		copy(sortedNames, names)
		sort.Strings(sortedNames)

		// Check for out-of-order variables
		for i, name := range names {
			if name != sortedNames[i] {
				var expectedAfter string
				if i > 0 {
					expectedAfter = sortedNames[i-1]
				}

				message := fmt.Sprintf("Variable '%s' is not in alphabetical order", name)
				if expectedAfter != "" {
					message = fmt.Sprintf("Variable '%s' is not in alphabetical order. Expected position after '%s'", name, expectedAfter)
				}

				err := EmitIssue(runner, r, message, blocks[i].DefRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
