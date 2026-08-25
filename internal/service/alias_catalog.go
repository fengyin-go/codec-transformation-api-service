package service

import "codec/internal/config"

type AliasCatalog struct {
	aliases   map[string]string
	validator config.AliasValidator
}

func NewAliasCatalog() *AliasCatalog {
	aliases, validator := config.LoadAliasDefaults()
	return &AliasCatalog{aliases: aliases, validator: validator}
}

func (c *AliasCatalog) Set(name, target string) error {
	if c.validator != nil {
		if err := c.validator.Validate(target); err != nil {
			return err
		}
	}
	c.aliases[name] = target
	return nil
}

func (c *AliasCatalog) Get(name string) string { return c.aliases[name] }
