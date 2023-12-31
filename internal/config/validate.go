package config

import (
	"fmt"
	"main/internal/core"
)

func (modules Modules) Validate(types map[string]core.Type) error {
	if err := validateModules(modules.Remote, types); err != nil {
		return err
	}
	return validateModules(modules.Local, types)
}

func validateModules(modules map[string]*string, types map[string]core.Type) error {
	for key, value := range modules {
		if value == nil {
			continue
		}
		if _, ok := types[*value]; !ok {
			return fmt.Errorf("incorrect type '%s' in module '%s'", *value, key)
		}
	}
	return nil
}
