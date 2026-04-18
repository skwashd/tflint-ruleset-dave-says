package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveIamNoInlinePolicyRule flags use of aws_iam_role_policy, aws_iam_user_policy,
// and aws_iam_group_policy. These create inline policies that are not reusable
// and not visible in the IAM console's managed policies view. Use aws_iam_policy
// with aws_iam_role_policy_attachment instead.
type DaveIamNoInlinePolicyRule struct {
	BaseRule
}

func NewDaveIamNoInlinePolicyRule() *DaveIamNoInlinePolicyRule {
	return &DaveIamNoInlinePolicyRule{
		BaseRule: BaseRule{ruleName: "dave_iam_no_inline_policy"},
	}
}

func (r *DaveIamNoInlinePolicyRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_iam_no_inline_policy.md"
}

var inlinePolicyResources = map[string]string{
	"aws_iam_role_policy":  "aws_iam_policy + aws_iam_role_policy_attachment",
	"aws_iam_user_policy":  "aws_iam_policy + aws_iam_user_policy_attachment",
	"aws_iam_group_policy": "aws_iam_policy + aws_iam_group_policy_attachment",
}

func (r *DaveIamNoInlinePolicyRule) Check(runner tflint.Runner) error {
	for resourceType, replacement := range inlinePolicyResources {
		resources, err := runner.GetResourceContent(resourceType, &hclext.BodySchema{}, nil)
		if err != nil {
			return err
		}

		for _, resource := range resources.Blocks {
			if err := runner.EmitIssue(r,
				"Use "+replacement+" instead of "+resourceType+". Inline policies are not reusable and are harder to audit in the IAM console.",
				resource.DefRange,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
