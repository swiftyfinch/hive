package check

import (
	"fmt"
	"hive/packages/common"
)

type validationFailure struct {
	ModuleName     string
	ModuleType     string
	DependencyName string
	DependencyType string
	IsWarning      bool
}

func checkDependencies(
	modules map[string]common.Module,
	bans []map[string]string,
	moduleTypes map[string]string,
) ([]validationFailure, error) {
	failures := []validationFailure{}
	for _, module := range modules {
		moduleFailures, err := checkModuleDependencies(
			module,
			bans,
			moduleTypes,
		)
		if err != nil {
			return nil, err
		}
		failures = append(failures, moduleFailures...)
	}
	return failures, nil
}

func checkModuleDependencies(
	module common.Module,
	bans []map[string]string,
	moduleTypes map[string]string,
) ([]validationFailure, error) {
	failures := []validationFailure{}
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
	return failures, nil
}
