package main

import (
	"github.com/skwashd/tflint-ruleset-dave-says/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &rules.DaveSaysRuleSet{
			BuiltinRuleSet: tflint.BuiltinRuleSet{
				Name:    "dave-says",
				Version: "0.2.5",
				Rules: []tflint.Rule{
					rules.NewDaveAwsPolicyNoJsonencodeRule(),
					rules.NewDaveCloudwatchLogRetentionRule(),
					rules.NewDaveIamNoInlinePolicyRule(),
					rules.NewDaveLabelMinLengthRule(),
					rules.NewDaveLabelNoTypeSubstringRule(),
					rules.NewDaveLabelSnakeRule(),
					rules.NewDaveListAlphabeticalOrderRule(),
					rules.NewDaveNoVpcIdVariableRule(),
					rules.NewDaveOutputMustBeInOutputsFileRule(),
					rules.NewDaveResourceNameKebabRule(),
					rules.NewDaveResourceNameNoTypeSubstringRule(),
					rules.NewDaveS3BucketNamespaceRule(),
					rules.NewDaveS3NoInlineConfigRule(),
					rules.NewDaveS3NoPublicAclRule(),
					rules.NewDaveSecurityGroupNoInlineRulesRule(),
					rules.NewDaveVariableAlphabeticalOrderRule(),
					rules.NewDaveVariableHasDescriptionRule(),
					rules.NewDaveVariableHasTypeRule(),
					rules.NewDaveVariableMustBeInVariablesFileRule(),
					rules.NewDaveVariableRegionRule(),
				},
			},
		},
	})
}
