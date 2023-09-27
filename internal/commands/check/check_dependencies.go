package check

import (
	"fmt"
	"main/internal/config"
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
	ignores []config.Ignore,
) ([]validationFailure, error) {
	failures := []validationFailure{}
	for _, module := range modules {
		moduleFailures, err := checkModuleDependencies(
			module,
			rules,
			moduleTypes,
			ignores,
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
	ignores []config.Ignore,
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

		isIgnored := false
		for _, ignore := range ignores {
			moduleRegexp, dependenciesRegexp := ignore.ModuleRegexp, ignore.DependenciesRegexp
			if moduleRegexp.MatchString(module.Name) && dependenciesRegexp.MatchString(dependency) {
				isIgnored = true
				break
			}
		}
		if isIgnored {
			// Validation of this module and dependency combination should be skipped
			continue
		}

		if !ok || !slices.Contains(allowDependencies, dependencyType) {
			failures = append(failures, validationFailure{
				module.Name, moduleType,
				dependency, dependencyType,
			})
		}
	}
	return failures, nil
}
