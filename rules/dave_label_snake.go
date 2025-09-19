package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveLabelSnakeRule struct {
	BaseRule
}

func NewDaveLabelSnakeRule() *DaveLabelSnakeRule {
	return &DaveLabelSnakeRule{
		BaseRule: BaseRule{ruleName: "dave_label_snake"},
	}
}

func (r *DaveLabelSnakeRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "data", LabelNames: []string{"type", "name"}},
			{Type: "ephemeral", LabelNames: []string{"type", "name"}},
			{Type: "module", LabelNames: []string{"name"}},
			{Type: "output", LabelNames: []string{"name"}},
			{Type: "resource", LabelNames: []string{"type", "name"}},
			{Type: "variable", LabelNames: []string{"name"}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) >= 2 {
			label := block.Labels[1]
			if !SnakeRegex.MatchString(label) {
				blockType := strings.ToUpper(block.Type[:1]) + block.Type[1:]
				err := EmitIssue(runner, r, fmt.Sprintf("%s label '%s' must contain only lowercase letters, numbers, and underscores", blockType, label), block.DefRange)
				if err != nil {
					return err
				}
			}
		}

		if len(block.Labels) == 1 {
			label := block.Labels[0]
			if !SnakeRegex.MatchString(label) {
				blockType := strings.ToUpper(block.Type[:1]) + block.Type[1:]
				err := EmitIssue(runner, r, fmt.Sprintf("%s name '%s' must contain only lowercase letters, numbers, and underscores", blockType, label), block.DefRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
