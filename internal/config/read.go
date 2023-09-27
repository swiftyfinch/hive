package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

func Read(path string) (*Config, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err = yaml.Unmarshal(buffer, config); err != nil {
		return nil, err
	}
	if config.Modules.Local == nil {
		config.Modules.Local = map[string]*string{}
	}
	if config.Modules.Remote == nil {
		config.Modules.Remote = map[string]*string{}
	}
	return config, err
}
