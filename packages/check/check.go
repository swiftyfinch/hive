package check

import (
	"fmt"
	"hive/packages/cocoapods"
	"hive/packages/config"
)

func Check(configPath string) error {
	config, err := config.Read(configPath)
	if err != nil {
		return err
	}
	if err := config.Validate(); err != nil {
		return err
	}

	// Read pods from Podfile.lock
	_, localPods, err := cocoapods.ReadPods()
	if err != nil {
		return err
	}

	// Merge modules from config
	moduleTypes, err := config.AllModulesTypes()
	if err != nil {
		return err
	}

	// Validate
	failures, err := checkDependencies(localPods, config.Bans, moduleTypes)
	if err != nil {
		return err
	}
	for _, failure := range failures {
		fmt.Println(formatMessage(failure))
	}

	return nil
}

type validationFailure struct {
	ModuleName     string
	ModuleType     string
	DependencyName string
	DependencyType string
	IsWarning      bool
}

func checkDependencies(
	modules map[string]cocoapods.Pod,
	bans []map[string]string,
	moduleTypes map[string]string,
) ([]validationFailure, error) {
	failures := []validationFailure{}
	for _, module := range modules {
		moduleType, ok := moduleTypes[module.Name]
		if !ok {
			return nil, fmt.Errorf("can't find type of module '%s'", module.Name)
		}
		for _, dependency := range module.Dependencies {
			dependencyType, ok := moduleTypes[dependency]
			if !ok {
				return nil, fmt.Errorf("can't find type of module '%s'", dependency)
			}
			for _, rule := range bans {
				if banDependency, ok := rule[moduleType]; ok && banDependency == dependencyType {
					severity, ok := rule["severity"]
					isWarning := ok && severity == "warning"
					failures = append(failures, validationFailure{
						module.Name, moduleType,
						dependency, dependencyType,
						isWarning,
					})
				}
			}
		}
	}
	return failures, nil
}

func formatMessage(failure validationFailure) string {
	moduleName, moduleType := failure.ModuleName, failure.ModuleType
	dependencyName, dependencyType := failure.DependencyName, failure.DependencyType
	return fmt.Sprintf("forbidden dependency %s(%s) → %s(%s)", moduleName, moduleType, dependencyName, dependencyType)
}
