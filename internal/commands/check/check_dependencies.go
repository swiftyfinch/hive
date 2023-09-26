package check

import (
	"fmt"
	"main/internal/core"
	"slices"
)

type validationFailure struct {
	ModuleName     string
	ModuleType     string
	DependencyName string
	DependencyType string
}

func checkDependencies(
	modules map[string]core.Module,
	rules map[string][]string,
	moduleTypes map[string]string,
) ([]validationFailure, error) {
	failures := []validationFailure{}
	for _, module := range modules {
		moduleFailures, err := checkModuleDependencies(
			module,
			rules,
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
	module core.Module,
	rules map[string][]string,
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

		allowDependencies, ok := rules[moduleType]
		if !ok {
			return nil, fmt.Errorf("unknown module type '%s'", moduleType)
		}

		if ok && slices.Contains(allowDependencies, dependencyType) {
			// Correct dependency
		} else {
			failures = append(failures, validationFailure{
				module.Name, moduleType,
				dependency, dependencyType,
			})
		}
	}
	return failures, nil
}
