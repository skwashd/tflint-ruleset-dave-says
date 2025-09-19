package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveLabelMinLengthRule struct {
	BaseRule
}

func NewDaveLabelMinLengthRule() *DaveLabelMinLengthRule {
	return &DaveLabelMinLengthRule{
		BaseRule: BaseRule{ruleName: "dave_label_min_length"},
	}
}

func (r *DaveLabelMinLengthRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "data", LabelNames: []string{"type", "name"}},
			{Type: "ephemeral", LabelNames: []string{"type", "name"}},
			{Type: "resource", LabelNames: []string{"type", "name"}}},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) >= 2 {
			label := block.Labels[1]
			if len(label) < 3 {
				if label == "db" || label == "s3" {
					continue
				}

				blockType := strings.ToUpper(block.Type[:1]) + block.Type[1:]
				err := EmitIssue(runner, r, fmt.Sprintf("%s label '%s' must be at least 3 characters long or a well known name such as db or s3", blockType, label), block.DefRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
