package tidy

import (
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/core"
	"os"
	"path/filepath"
)

func Tidy(workingDirectory string) error {
	configPath := workingDirectory + "/" + core.Modules_File_Name
	config, err := getConfig(configPath)
	if err != nil {
		return err
	}

	types := core.DefaultTypes()
	if err := config.Validate(types); err != nil {
		return err
	}

	// Read modules from Podfile.lock
	remoteModules, localModules, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from Podfile.lock and config
	updateModules(remoteModules, &config.Modules.Remote, types)
	updateModules(localModules, &config.Modules.Local, types)

	// Save updated config
	return writeConfig(*config, configPath)
}

func getConfig(path string) (*config.Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &config.Config{
			Modules: config.Modules{
				Remote: map[string]*string{},
				Local:  map[string]*string{},
			},
		}, nil
	}
	return config.Read(path)
}

func updateModules(
	modules map[string]core.Module,
	configModules *map[string]*string,
	types map[string]core.Type,
) error {
	// Remove redundant modules
	for module := range *configModules {
		if _, ok := modules[module]; !ok {
			delete(*configModules, module)
		}
	}

	for _, module := range modules {
		match := false
		for _, moduleType := range types {
			// Try to find out what type of module is it
			for _, regexp := range moduleType.Regexps {
				if regexp.MatchString(module.Name) {
					(*configModules)[module.Name] = &moduleType.Name
					match = true
					break
				}
			}
			if match {
				break
			}
		}

		if _, ok := (*configModules)[module.Name]; !ok && !match {
			// Add new module with null type
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
