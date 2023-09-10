package config

import (
	"bytes"
	"os"

	"gopkg.in/yaml.v3"
)

func WriteModules(modules Modules, path string) error {
	return writeYMLFile(modules, path)
}

type Ban struct {
	ModuleType     string
	DependencyType string
	Severity       string
}
type rulesYML struct {
	Types []string            `yaml:"types"`
	Bans  []map[string]string `yaml:"bans"`
}

func WriteRules(types []string, bans []Ban, path string) error {
	bansMap := []map[string]string{}
	for _, ban := range bans {
		bansMap = append(bansMap, map[string]string{
			ban.ModuleType: ban.DependencyType,
			"severity":     ban.Severity,
		})
	}
	return writeYMLFile(
		rulesYML{types, bansMap},
		path,
	)
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
