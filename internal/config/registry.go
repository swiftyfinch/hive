package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadRegistry(path string) (map[string]string, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	modules := map[string]string{}
	if err = yaml.Unmarshal(buffer, modules); err != nil {
		return nil, err
	}
	return modules, err
}
