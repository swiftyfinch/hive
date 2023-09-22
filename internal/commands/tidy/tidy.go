package tidy

import (
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/modules"
	"os"
	"path/filepath"
	"regexp"
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
	updateModules(remoteModules, &config.Modules.Remote, config.Types)
	updateModules(localModules, &config.Modules.Local, config.Types)

	// Save updated config
	return writeConfig(*config, configPath)
}

func readConfig(path string) (*config.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		types := []interface{}{}
		for _, value := range []string{"base", "feature"} {
			types = append(types, value)
		}

		return &config.Config{
			Types: types,
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
	modules map[string]modules.Module,
	configModules *map[string]*string,
	types []interface{},
) error {
	for module := range *configModules {
		if _, ok := modules[module]; !ok {
			delete(*configModules, module)
		}
	}

	for _, module := range modules {
		for _, element := range types {
			if regex := config.TypeRegex(element); regex != nil {
				regex, err := regexp.Compile(*regex)
				if err != nil {
					return err
				}
				if regex.MatchString(module.Name) {
					(*configModules)[module.Name] = config.TypeValue(element)
					break
				}
			}
		}

		if _, ok := (*configModules)[module.Name]; !ok {
			(*configModules)[module.Name] = nil
		}
	}
	return nil
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
