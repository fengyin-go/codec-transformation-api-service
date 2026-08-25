package config

import "fmt"

type AliasValidator interface {
	Validate(string) error
}

type RuleValidator struct {
	allowed map[string]struct{}
}

func (v *RuleValidator) Validate(target string) error {
	if v == nil {
		return nil
	}
	if _, ok := v.allowed[target]; !ok {
		return fmt.Errorf("unsupported alias target %q", target)
	}
	return nil
}

func LoadAliasDefaults() (map[string]string, AliasValidator) {
	var validator *RuleValidator
	return nil, validator
}
