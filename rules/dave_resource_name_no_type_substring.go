package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveResourceNameNoTypeSubstringRule struct {
	BaseRule
}

func NewDaveResourceNameNoTypeSubstringRule() *DaveResourceNameNoTypeSubstringRule {
	return &DaveResourceNameNoTypeSubstringRule{
		BaseRule: BaseRule{ruleName: "dave_resource_name_no_type_substring"},
	}
}

func (r *DaveResourceNameNoTypeSubstringRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_resource_name_no_type_substring.md"
}

func (r *DaveResourceNameNoTypeSubstringRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					{Name: "name"},
					{Name: "name_prefix"},
				},
			}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) >= 2 {
			resourceType := block.Labels[0]
			typeWords := SplitWords(resourceType)

			// Check name attribute
			if nameAttr, exists := block.Body.Attributes["name"]; exists {
				if err := r.checkNameAttribute(runner, resourceType, typeWords, "name", nameAttr); err != nil {
					return err
				}
			}

			// Check name_prefix attribute
			if namePrefixAttr, exists := block.Body.Attributes["name_prefix"]; exists {
				if err := r.checkNameAttribute(runner, resourceType, typeWords, "name_prefix", namePrefixAttr); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *DaveResourceNameNoTypeSubstringRule) checkNameAttribute(runner tflint.Runner, resourceType string, typeWords []string, attrName string, attr *hclext.Attribute) error {
	var nameValue string
	err := runner.EvaluateExpr(attr.Expr, &nameValue, nil)
	if err != nil {
		// Skip if we can't evaluate (variables, references, etc.)
		return nil
	}

	nameWords := SplitWordsOnDash(nameValue)  // Names only have dashes
	if found, word := ContainsAnyWord(nameWords, typeWords); found {
		return EmitIssue(runner, r, fmt.Sprintf("Resource %s attribute '%s' contains substring '%s' from resource type '%s'", attrName, nameValue, word, resourceType), attr.Range)
	}

	return nil
}
