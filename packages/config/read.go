package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ReadConfig(path string) (*Config, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = yaml.Unmarshal(buffer, config); err != nil {
		return nil, err
	}
	return config, err
}
