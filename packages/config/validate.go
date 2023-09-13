package config

import (
	"fmt"
)

func (config Config) Validate() error {
	if err := config.validateRules(); err != nil {
		return err
	}
	return config.validateModules()
}

func typesContains(types []interface{}, value string) bool {
	for _, element := range types {
		switch resolvedType := element.(type) {
		case string:
			if element == value {
				return true
			}
		case map[string]interface{}:
			if _, ok := resolvedType[value]; ok {
				return true
			}
		}
	}
	return false
}

func (config Config) validateRules() error {
	for _, ban := range config.Bans {
		for key, value := range ban {
			if key == "severity" {
				if value != "error" && value != "warning" {
					return fmt.Errorf("incorrect severity '%s'", value)
				}
			} else {
				if !typesContains(config.Types, key) {
					return fmt.Errorf("incorrect module type in rule '%s'", key)
				}
				if !typesContains(config.Types, value) {
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

func validateModules(modules map[string]*string, types []interface{}) error {
	for key, value := range modules {
		if value != nil && !typesContains(types, *value) {
			return fmt.Errorf("incorrect type '%s' in module '%s'", *value, key)
		}
	}
	return nil
}
