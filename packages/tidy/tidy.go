package tidy

import (
	"fmt"
	"hive/packages/cocoapods"
	"hive/packages/config"
	glob "hive/packages/utils"
	"os"
	"path/filepath"
)

const basePath = ".devtools/hive"
const modulesPath = basePath + "/modules"
const localModulesPath = modulesPath + "/local.yml"
const remoteModulesPath = modulesPath + "/remote.yml"
const rulesPath = basePath + "/rules.yml"

var types = []string{"base", "feature"}
var bans = []config.Ban{
	{ModuleType: "feature", DependencyType: "feature", Severity: "error"},
	{ModuleType: "base", DependencyType: "feature", Severity: "error"},
	{ModuleType: "base", DependencyType: "base", Severity: "warning"},
}

func Tidy() error {
	remotePods, localPods, err := readPods()
	if err != nil {
		return err
	}

	if err = updateModules(remotePods, remoteModulesPath); err != nil {
		return err
	}
	if err = updateModules(localPods, localModulesPath); err != nil {
		return err
	}
	if err = writeRules(types, bans, rulesPath); err != nil {
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

func updateModules(
	pods map[string]cocoapods.Pod,
	path string,
) error {
	directoryPath := filepath.Dir(path)
	if _, err := os.Stat(directoryPath); os.IsNotExist(err) {
		if err := os.MkdirAll(directoryPath, os.ModePerm); err != nil {
			return err
		}
	}

	modules, err := config.ReadModules(path)
	if err != nil {
		return err
	}
	for _, pod := range pods {
		if _, ok := modules[pod.Name]; !ok {
			modules[pod.Name] = nil
		}
	}
	return config.WriteModules(modules, path)
}

func writeRules(types []string, bans []config.Ban, path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return config.WriteRules(types, bans, path)
	}
	return nil
}
