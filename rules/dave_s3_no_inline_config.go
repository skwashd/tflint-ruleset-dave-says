package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// DaveS3NoInlineConfigRule flags deprecated inline configuration blocks on aws_s3_bucket.
type DaveS3NoInlineConfigRule struct {
	BaseRule
}

func NewDaveS3NoInlineConfigRule() *DaveS3NoInlineConfigRule {
	return &DaveS3NoInlineConfigRule{
		BaseRule: BaseRule{ruleName: "dave_s3_no_inline_config"},
	}
}

func (r *DaveS3NoInlineConfigRule) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *DaveS3NoInlineConfigRule) Link() string {
	return "https://github.com/skwashd/tflint-ruleset-dave-says/blob/main/docs/rules/dave_s3_no_inline_config.md"
}

// deprecatedS3Blocks lists inline blocks on aws_s3_bucket that were deprecated
// in AWS provider v4 and should use dedicated resources instead.
var deprecatedS3Blocks = map[string]string{
	"versioning":                           "aws_s3_bucket_versioning",
	"logging":                              "aws_s3_bucket_logging",
	"server_side_encryption_configuration": "aws_s3_bucket_server_side_encryption_configuration",
	"lifecycle_rule":                       "aws_s3_bucket_lifecycle_configuration",
	"cors_rule":                            "aws_s3_bucket_cors_configuration",
	"website":                              "aws_s3_bucket_website_configuration",
	"replication_configuration":            "aws_s3_bucket_replication_configuration",
	"object_lock_configuration":            "aws_s3_bucket_object_lock_configuration",
	"grant":                                "aws_s3_bucket_acl",
}

var deprecatedS3Attrs = map[string]string{
	"acl":                 "aws_s3_bucket_acl",
	"policy":              "aws_s3_bucket_policy",
	"acceleration_status": "aws_s3_bucket_accelerate_configuration",
	"request_payer":       "aws_s3_bucket_request_payment_configuration",
}

func (r *DaveS3NoInlineConfigRule) Check(runner tflint.Runner) error {
	blockSchemas := make([]hclext.BlockSchema, 0, len(deprecatedS3Blocks))
	for blockType := range deprecatedS3Blocks {
		blockSchemas = append(blockSchemas, hclext.BlockSchema{
			Type: blockType,
		})
	}

	attrSchemas := make([]hclext.AttributeSchema, 0, len(deprecatedS3Attrs))
	for attrName := range deprecatedS3Attrs {
		attrSchemas = append(attrSchemas, hclext.AttributeSchema{Name: attrName})
	}

	resources, err := runner.GetResourceContent("aws_s3_bucket", &hclext.BodySchema{
		Attributes: attrSchemas,
		Blocks:     blockSchemas,
	}, nil)
	if err != nil {
		return err
	}

	for _, resource := range resources.Blocks {
		for _, block := range resource.Body.Blocks {
			replacement, ok := deprecatedS3Blocks[block.Type]
			if !ok {
				continue
			}
			if err := runner.EmitIssue(r,
				fmt.Sprintf("Inline %q block is deprecated since AWS provider v4. Use the %s resource instead.", block.Type, replacement),
				block.DefRange,
			); err != nil {
				return err
			}
		}

		for attrName, replacement := range deprecatedS3Attrs {
			attr, exists := resource.Body.Attributes[attrName]
			if !exists {
				continue
			}
			if err := runner.EmitIssue(r,
				fmt.Sprintf("Inline %q attribute is deprecated since AWS provider v4. Use the %s resource instead.", attrName, replacement),
				attr.Range,
			); err != nil {
				return err
			}
		}
	}
	return nil
}
