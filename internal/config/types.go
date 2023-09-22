package config

import "fmt"

type ModulesMap map[string]*string

type Modules struct {
	Remote map[string]*string `yaml:"remote"`
	Local  map[string]*string `yaml:"local"`
}

type Config struct {
	Types   []interface{}       `yaml:"types"`
	Bans    []map[string]string `yaml:"bans"`
	Modules Modules             `yaml:"modules"`
}

func TypeValue(theType interface{}) *string {
	switch resolvedType := theType.(type) {
	case string:
		return &resolvedType
	case map[string]interface{}:
		for key := range resolvedType {
			return &key
		}
	}
	return nil
}
func TypeRegex(theType interface{}) *string {
	switch resolvedType := theType.(type) {
	case string:
		return nil
	case map[string]interface{}:
		for _, value := range resolvedType {
			if stringValue, ok := value.(string); ok {
				return &stringValue
			}
		}
	}
	return nil
}

func (config Config) AllModulesTypes() (map[string]string, error) {
	moduleTypes := map[string]string{}
	for key, value := range config.Modules.Remote {
		if value == nil {
			return nil, fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}
	for key, value := range config.Modules.Local {
		if value == nil {
			return nil, fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}
	return moduleTypes, nil
}
