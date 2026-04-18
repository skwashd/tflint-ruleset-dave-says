package rules

import (
	"fmt"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// ConfigurableRule is implemented by rules that accept per-rule config
// from .tflint.hcl rule blocks.
type ConfigurableRule interface {
	ApplyRuleConfig(body *hclext.BodyContent) error
}

// DaveSaysRuleSet extends BuiltinRuleSet to support per-rule configuration.
type DaveSaysRuleSet struct {
	tflint.BuiltinRuleSet
}

func (r *DaveSaysRuleSet) ApplyConfig(content *hclext.BodyContent) error {
	// Apply standard config (enabled/disabled per rule)
	if err := r.BuiltinRuleSet.ApplyConfig(content); err != nil {
		return err
	}

	// Apply rule-specific config to configurable rules
	for _, rule := range r.EnabledRules {
		configurable, ok := rule.(ConfigurableRule)
		if !ok {
			continue
		}

		for _, block := range content.Blocks {
			if len(block.Labels) > 0 && block.Labels[0] == rule.Name() {
				if err := configurable.ApplyRuleConfig(block.Body); err != nil {
					return fmt.Errorf("configuring rule %s: %w", rule.Name(), err)
				}
			}
		}
	}
	return nil
}
