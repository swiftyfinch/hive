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

	// Read pods from Podfile.lock
	remotePods, localPods, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from Podfile.lock and config
	updateModules(remotePods, &config.Modules.Remote)
	updateModules(localPods, &config.Modules.Local)

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
	pods map[string]common.Module,
	modules *map[string]*string,
) {
	for module := range *modules {
		if _, ok := pods[module]; !ok {
			delete(*modules, module)
		}
	}
	for _, pod := range pods {
		if _, ok := (*modules)[pod.Name]; !ok {
			(*modules)[pod.Name] = nil
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
