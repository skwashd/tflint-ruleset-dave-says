package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveS3NoPublicAclRule flags S3 bucket ACLs that allow public access.
// When run with --fix, replaces the public ACL with "private".
type DaveS3NoPublicAclRule struct {
	BaseRule
}

func NewDaveS3NoPublicAclRule() *DaveS3NoPublicAclRule {
	return &DaveS3NoPublicAclRule{
		BaseRule: BaseRule{ruleName: "dave_s3_no_public_acl"},
	}
}

func (r *DaveS3NoPublicAclRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *DaveS3NoPublicAclRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_s3_no_public_acl.md"
}

var publicACLs = map[string]bool{
	"public-read":       true,
	"public-read-write": true,
	"authenticated-read": true,
}

func (r *DaveS3NoPublicAclRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_s3_bucket_acl", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "acl"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attr, exists := resource.Body.Attributes["acl"]
		if !exists {
			continue
		}

		var acl string
		if err := runner.EvaluateExpr(attr.Expr, &acl, nil); err != nil {
			continue
		}

		if publicACLs[acl] {
			if err := runner.EmitIssueWithFix(r,
				fmt.Sprintf("S3 bucket ACL is %q. Public buckets are a leading cause of data breaches. Use CloudFront for web access.", acl),
				attr.Expr.Range(),
				func(f tflint.Fixer) error {
					return f.ReplaceText(attr.Expr.Range(), `"private"`)
				},
			); err != nil {
				return err
			}
		}
	}
	return nil
}
