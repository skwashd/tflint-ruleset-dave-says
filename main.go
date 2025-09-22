package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/yourusername/tflint-ruleset-dave-says/rules"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "dave-says",
			Version: "0.1.0",
			Rules: []tflint.Rule{
				rules.NewDaveAwsPolicyNoJsonencodeRule(),
				rules.NewDaveLabelMinLengthRule(),
				rules.NewDaveLabelNoTypeSubstringRule(),
				rules.NewDaveLabelSnakeRule(),
				rules.NewDaveResourceNameKebabRule(),
				rules.NewDaveResourceNameNoTypeSubstringRule(),
				rules.NewDaveVariableAlphabeticalOrderRule(),
				rules.NewDaveVariableMustBeInVariablesFileRule(),
				rules.NewDaveVariableRegionRule(),
			},
		},
	})
}
