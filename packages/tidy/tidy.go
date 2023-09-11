package tidy

import (
	"fmt"
	"hive/packages/cocoapods"
	"hive/packages/config"
	glob "hive/packages/utils"
	"os"
	"path/filepath"
)

const configPath = ".devtools/hive.yml"

func Tidy() error {
	remotePods, localPods, err := readPods()
	if err != nil {
		return err
	}

	var config config.Config
	configPtr, err := readConfig(configPath)
	if err != nil {
		return err
	}
	config = *configPtr

	updateModules(remotePods, &config.Modules.Remote)
	updateModules(localPods, &config.Modules.Local)

	if err := writeConfig(config, configPath); err != nil {
		return err
	}

	return nil
}

func readPods() (
	remotePods map[string]cocoapods.Pod,
	localPods map[string]cocoapods.Pod,
	err error,
) {
	paths, err := glob.FindPathsRecursively(".", "Podfile.lock")
	if err != nil {
		return nil, nil, err
	}
	if len(paths) == 0 {
		return nil, nil, fmt.Errorf("couldn't find any Podfile.lock")
	} else if len(paths) > 1 {
		return nil, nil, fmt.Errorf("found several Podfile.lock files")
	}
	return cocoapods.ParsePodfile(paths[0])
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
	return config.ReadConfig(path)
}

func updateModules(
	pods map[string]cocoapods.Pod,
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
	config.WriteConfig(path)
	return nil
}
