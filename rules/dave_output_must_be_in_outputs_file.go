package rules

import (
	"fmt"
	"path/filepath"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveOutputMustBeInOutputsFileRule ensures all output blocks are declared in outputs.tf.
type DaveOutputMustBeInOutputsFileRule struct {
	BaseRule
}

func NewDaveOutputMustBeInOutputsFileRule() *DaveOutputMustBeInOutputsFileRule {
	return &DaveOutputMustBeInOutputsFileRule{
		BaseRule: BaseRule{ruleName: "dave_output_must_be_in_outputs_file"},
	}
}

func (r *DaveOutputMustBeInOutputsFileRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_output_must_be_in_outputs_file.md"
}

func (r *DaveOutputMustBeInOutputsFileRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{
				Type:       "output",
				LabelNames: []string{"name"},
				Body:       &hclext.BodySchema{},
			},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		filename := filepath.Base(block.DefRange.Filename)
		if filename != "outputs.tf" {
			if err := EmitIssue(runner, r,
				fmt.Sprintf("Output %q is defined in %s. All outputs should be in outputs.tf for consistent file organization.", block.Labels[0], filename),
				block.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
