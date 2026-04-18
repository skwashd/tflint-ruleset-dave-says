package rules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type DaveAwsPolicyNoJsonencodeRule struct {
	BaseRule
}

func NewDaveAwsPolicyNoJsonencodeRule() *DaveAwsPolicyNoJsonencodeRule {
	return &DaveAwsPolicyNoJsonencodeRule{
		BaseRule: BaseRule{ruleName: "dave_aws_policy_no_jsonencode"},
	}
}

func (r *DaveAwsPolicyNoJsonencodeRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_aws_policy_no_jsonencode.md"
}

func (r *DaveAwsPolicyNoJsonencodeRule) Check(runner tflint.Runner) error {
	content, err := runner.GetModuleContent(&hclext.BodySchema{
		Blocks: []hclext.BlockSchema{
			{Type: "resource", LabelNames: []string{"type", "name"}, Body: &hclext.BodySchema{
				Attributes: []hclext.AttributeSchema{{Name: "assume_role_policy"}, {Name: "policy"}, {Name: "bucket_policy"}},
			}},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, block := range content.Blocks {
		if len(block.Labels) >= 2 {
			resourceType := block.Labels[0]
			
			// Only check AWS resources
			if !strings.HasPrefix(resourceType, "aws_") {
				continue
			}

			// Check specific policy attributes
			for name, attr := range block.Body.Attributes {
				if r.isPolicyAttribute(name) {
					if r.containsJsonencode(attr.Expr) {
						err := EmitIssue(runner, r, fmt.Sprintf("AWS resource '%s' attribute '%s' should reference an aws_iam_policy_document data source instead of using jsonencode()", resourceType, name), attr.Range)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

func (r *DaveAwsPolicyNoJsonencodeRule) isPolicyAttribute(name string) bool {
	return name == "policy" || name == "assume_role_policy" || name == "bucket_policy" || strings.HasSuffix(name, "_policy")
}

func (r *DaveAwsPolicyNoJsonencodeRule) containsJsonencode(expr hcl.Expression) bool {
	switch e := expr.(type) {
	case *hclsyntax.FunctionCallExpr:
		return e.Name == "jsonencode"
	case *hclsyntax.ParenthesesExpr:
		return r.containsJsonencode(e.Expression)
	}
	return false
}
