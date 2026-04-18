package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveS3BucketNamespaceRule ensures S3 buckets use account-regional namespace.
type DaveS3BucketNamespaceRule struct {
	BaseRule
}

func NewDaveS3BucketNamespaceRule() *DaveS3BucketNamespaceRule {
	return &DaveS3BucketNamespaceRule{
		BaseRule: BaseRule{ruleName: "dave_s3_bucket_namespace"},
	}
}

func (r *DaveS3BucketNamespaceRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_s3_bucket_namespace.md"
}

func (r *DaveS3BucketNamespaceRule) Check(runner tflint.Runner) error {
	resources, err := runner.GetResourceContent("aws_s3_bucket", &hclext.BodySchema{
		Attributes: []hclext.AttributeSchema{
			{Name: "bucket_namespace"},
		},
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		attr, exists := resource.Body.Attributes["bucket_namespace"]
		if !exists {
			if err := runner.EmitIssue(r,
				`S3 buckets should use bucket_namespace = "account-regional" to avoid global name collisions and confused-deputy risks. Requires AWS provider >= 6.37.0.`,
				resource.DefRange,
			); err != nil {
				return err
			}
			continue
		}

		var namespace string
		if err := runner.EvaluateExpr(attr.Expr, &namespace, nil); err != nil {
			continue
		}

		if namespace != "account-regional" {
			if err := runner.EmitIssue(r,
				`S3 bucket_namespace should be "account-regional", not "`+namespace+`"`,
				attr.Expr.Range(),
			); err != nil {
				return err
			}
		}
	}
	return nil
}
