package config

import (
	"bytes"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Modules struct {
	Remote map[string]*string `yaml:"remote"`
	Local  map[string]*string `yaml:"local"`
}

func ReadModules(path string) (*Modules, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modules := &Modules{}
	if err = yaml.Unmarshal(buffer, modules); err != nil {
		return nil, err
	}
	if modules.Local == nil {
		modules.Local = map[string]*string{}
	}
	if modules.Remote == nil {
		modules.Remote = map[string]*string{}
	}
	return modules, err
}

func (modules Modules) Write(path string) error {
	var buffer bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buffer)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&modules)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (modules Modules) Types() (map[string]string, error) {
	moduleTypes := map[string]string{}
	for key, value := range modules.Remote {
		if value == nil {
			return nil, fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}
	for key, value := range modules.Local {
		if value == nil {
			return nil, fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}
	return moduleTypes, nil
}
