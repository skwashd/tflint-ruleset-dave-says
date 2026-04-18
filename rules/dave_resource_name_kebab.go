package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveResourceNameKebabRule struct {
	BaseRule
}

func NewDaveResourceNameKebabRule() *DaveResourceNameKebabRule {
	return &DaveResourceNameKebabRule{
		BaseRule: BaseRule{ruleName: "dave_resource_name_kebab"},
	}
}

func (r *DaveResourceNameKebabRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_resource_name_kebab.md"
}

func (r *DaveResourceNameKebabRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{
					// TODO: Identify other common name attributes
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
			// Check name attribute
			if nameAttr, exists := block.Body.Attributes["name"]; exists {
				if err := r.checkNameAttribute(runner, "name", nameAttr); err != nil {
					return err
				}
			}

			// Check name_prefix attribute
			if namePrefixAttr, exists := block.Body.Attributes["name_prefix"]; exists {
				if err := r.checkNameAttribute(runner, "name_prefix", namePrefixAttr); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (r *DaveResourceNameKebabRule) checkNameAttribute(runner tflint.Runner, attrName string, attr *hclext.Attribute) error {
	var nameValue string
	err := runner.EvaluateExpr(attr.Expr, &nameValue, nil)
	if err != nil {
		// Skip if we can't evaluate (variables, references, etc.)
		return nil
	}

	if !KebabRegex.MatchString(nameValue) {
		return EmitIssue(runner, r, fmt.Sprintf("Resource %s attribute '%s' must contain only lowercase letters, numbers, and dashes", attrName, nameValue), attr.Range)
	}

	return nil
}
