package config

import (
	"fmt"
	"main/internal/core"
)

func (config Config) Validate(types map[string]core.Type) error {
	return config.validateModules(types)
}

func (config Config) validateModules(types map[string]core.Type) error {
	if err := validateModules(config.Modules.Remote, types); err != nil {
		return err
	}
	return validateModules(config.Modules.Local, types)
}

func validateModules(modules map[string]*string, types map[string]core.Type) error {
	for key, value := range modules {
		if value == nil {
			return fmt.Errorf("nil type in module '%s'", key)
		}
		if _, ok := types[*value]; !ok {
			return fmt.Errorf("incorrect type '%s' in module '%s'", *value, key)
		}
	}
	return nil
}
