package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

func (config Config) Write(path string) error {
	return writeYMLFile(config, path)
}

func writeYMLFile(content interface{}, path string) error {
	var buffer bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&buffer)
	yamlEncoder.SetIndent(2)
	err := yamlEncoder.Encode(&content)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}
	return nil
}
