package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadModules(path string) (Modules, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modules := Modules{}
	if err = yaml.Unmarshal(buffer, modules); err != nil {
		return nil, err
	}
	return modules, err
}
