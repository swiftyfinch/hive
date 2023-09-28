package tidy

import (
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/core"
	"os"
	"path/filepath"
)

func Tidy(workingDirectory string) error {
	modulesPath := workingDirectory + "/" + core.Modules_File_Name
	modules, err := getModules(modulesPath)
	if err != nil {
		return err
	}

	types := core.DefaultTypes()
	if err := modules.Validate(types); err != nil {
		return err
	}

	// Read modules from Podfile.lock
	remoteModules, localModules, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from Podfile.lock and config
	updateModules(remoteModules, &modules.Remote, types)
	updateModules(localModules, &modules.Local, types)

	// Save updated config
	return writeModules(*modules, modulesPath)
}

func getModules(path string) (*config.Modules, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &config.Modules{
			Remote: map[string]*string{},
			Local:  map[string]*string{},
		}, nil
	}
	return config.ReadModules(path)
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

func writeModules(config config.Modules, path string) error {
	directoryPath := filepath.Dir(path)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
			return err
		}
	}
	config.Write(path)
	return nil
}
