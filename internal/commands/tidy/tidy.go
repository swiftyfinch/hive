package tidy

import (
	"fmt"
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/core"
	"os"
	"path/filepath"
)

func Tidy(workingDirectory string, registryPath *string) error {
	modulesPath := workingDirectory + "/" + core.Modules_File_Name
	cachedModules, err := getModules(modulesPath)
	if err != nil {
		return err
	}

	types := core.DefaultTypes()
	if err := cachedModules.Validate(types); err != nil {
		return err
	}

	// Read modules from Podfile.lock
	remoteModules, localModules, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from Podfile.lock and config
	updateModules(remoteModules, &cachedModules.Remote, types)
	updateModules(localModules, &cachedModules.Local, types)

	// Override cached modules from registy
	if registryPath != nil {
		registyModules, err := config.ReadRegistry(*registryPath)
		fmt.Println(registyModules)
		if err != nil {
			return err
		}
		updateRemoteModules(&cachedModules.Remote, registyModules)
	}

	// Save updated config
	return writeModules(*cachedModules, modulesPath)
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

func updateRemoteModules(
	cachedModules *map[string]*string,
	registyModules map[string]string,
) {
	for name := range *cachedModules {
		if moduleType, ok := registyModules[name]; ok {
			(*cachedModules)[name] = &moduleType
		}
	}
}

func updateModules(
	modules map[string]core.Module,
	cachedModules *map[string]*string,
	types map[string]core.Type,
) error {
	// Remove redundant modules
	for module := range *cachedModules {
		if _, ok := modules[module]; !ok {
			delete(*cachedModules, module)
		}
	}

	for _, module := range modules {
		match := false
		for _, moduleType := range types {
			// Try to find out what type of module is it
			for _, regexp := range moduleType.Regexps {
				if regexp.MatchString(module.Name) {
					(*cachedModules)[module.Name] = &moduleType.Name
					match = true
					break
				}
			}
			if match {
				break
			}
		}

		if _, ok := (*cachedModules)[module.Name]; !ok && !match {
			// Add new module with null type
			(*cachedModules)[module.Name] = nil
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
