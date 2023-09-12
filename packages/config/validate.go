package config

import (
	"fmt"

	"golang.org/x/exp/slices"
)

func (config Config) Validate() error {
	if err := config.validateRules(); err != nil {
		return err
	}
	return config.validateModules()
}

func (config Config) validateRules() error {
	for _, ban := range config.Bans {
		for key, value := range ban {
			if key == "severity" {
				if value != "error" && value != "warning" {
					return fmt.Errorf("incorrect severity '%s'", value)
				}
			} else {
				if !slices.Contains(config.Types, key) {
					return fmt.Errorf("incorrect module type in rule '%s'", key)
				}
				if !slices.Contains(config.Types, value) {
					return fmt.Errorf("incorrect dependency type in rule '%s'", value)
				}
			}
		}
	}
	return nil
}

func (config Config) validateModules() error {
	if err := validateModules(config.Modules.Remote, config.Types); err != nil {
		return err
	}
	return validateModules(config.Modules.Local, config.Types)
}

func validateModules(modules map[string]*string, types []string) error {
	for key, value := range modules {
		if value != nil && !slices.Contains(types, *value) {
			return fmt.Errorf("incorrect type '%s' in module '%s'", *value, key)
		}
	}
	return nil
}
