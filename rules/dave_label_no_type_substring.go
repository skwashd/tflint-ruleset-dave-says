package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveLabelNoTypeSubstringRule struct {
	BaseRule
}

func NewDaveLabelNoTypeSubstringRule() *DaveLabelNoTypeSubstringRule {
	return &DaveLabelNoTypeSubstringRule{
		BaseRule: BaseRule{ruleName: "dave_label_no_type_substring"},
	}
}

func (r *DaveLabelNoTypeSubstringRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_label_no_type_substring.md"
}

func (r *DaveLabelNoTypeSubstringRule) Check(runner tflint.Runner) error {
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
			resourceType := block.Labels[0]
			label := block.Labels[1]

			typeWords := SplitWords(resourceType)       // Resource types can have both _ and -
			labelWords := SplitWordsOnUnderscore(label) // Labels only have underscores

			if found, word := ContainsAnyWord(labelWords, typeWords); found {
				blockType := strings.ToUpper(block.Type[:1]) + block.Type[1:]
				err := EmitIssue(runner, r, fmt.Sprintf("%s label '%s' contains substring '%s' from %s type '%s'", blockType, label, word, strings.ToLower(blockType), resourceType), block.DefRange)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}
