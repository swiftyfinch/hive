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
	moduleTypes := map[string]string{}
	for key, value := range config.Modules.Remote {
		if value == nil {
			return fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}
	for key, value := range config.Modules.Local {
		if value == nil {
			return fmt.Errorf("module '%s' has empty type", key)
		}
		moduleTypes[key] = *value
	}

	// Validate
	failures := []ValidationFailure{}
	for _, module := range localPods {
		moduleType, ok := moduleTypes[module.Name]
		if !ok {
			return fmt.Errorf("can't find type of module '%s'", module.Name)
		}
		for _, dependency := range module.Dependencies {
			dependencyType, ok := moduleTypes[dependency]
			if !ok {
				return fmt.Errorf("can't find type of module '%s'", dependency)
			}
			for _, rule := range config.Bans {
				if banDependency, ok := rule[moduleType]; ok && banDependency == dependencyType {
					severity, ok := rule["severity"]
					isWarning := ok && severity == "warning"
					failures = append(failures, ValidationFailure{
						module.Name, moduleType,
						dependency, dependencyType,
						isWarning,
					})
				}
			}
		}
	}

	for _, failure := range failures {
		fmt.Println(formatMessage(failure))
	}

	return nil
}

type ValidationFailure struct {
	ModuleName     string
	ModuleType     string
	DependencyName string
	DependencyType string
	IsWarning      bool
}

func formatMessage(failure ValidationFailure) string {
	moduleName, moduleType := failure.ModuleName, failure.ModuleType
	dependencyName, dependencyType := failure.DependencyName, failure.DependencyType
	return fmt.Sprintf("forbidden dependency %s(%s) → %s(%s)", moduleName, moduleType, dependencyName, dependencyType)
}
