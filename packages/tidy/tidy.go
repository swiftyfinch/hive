package tidy

import (
	"hive/packages/cocoapods"
	"hive/packages/common"
	"hive/packages/config"
	"os"
	"path/filepath"
)

func Tidy(configPath string) error {
	config, err := readConfig(configPath)
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return err
	}

	// Read modules from Podfile.lock
	remoteModules, localModules, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from Podfile.lock and config
	updateModules(remoteModules, &config.Modules.Remote)
	updateModules(localModules, &config.Modules.Local)

	// Save updated config
	return writeConfig(*config, configPath)
}

func readConfig(path string) (*config.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &config.Config{
			Types: []string{"base", "feature"},
			Bans: []map[string]string{
				{"feature": "feature", "severity": "error"},
				{"base": "feature", "severity": "error"},
				{"base": "base", "severity": "warning"},
			},
			Modules: config.Modules{
				Remote: map[string]*string{},
				Local:  map[string]*string{},
			},
		}, nil
	}
	return config.Read(path)
}

func updateModules(
	modules map[string]common.Module,
	configModules *map[string]*string,
) {
	for module := range *configModules {
		if _, ok := modules[module]; !ok {
			delete(*configModules, module)
		}
	}
	for _, module := range modules {
		if _, ok := (*configModules)[module.Name]; !ok {
			(*configModules)[module.Name] = nil
		}
	}
}

func writeConfig(config config.Config, path string) error {
	directoryPath := filepath.Dir(path)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
			return err
		}
	}
	config.Write(path)
	return nil
}
