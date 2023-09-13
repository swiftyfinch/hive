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

	// Check
	failures, err := checkDependencies(localPods, config.Bans, moduleTypes)
	if err != nil {
		return err
	}
	for _, failure := range failures {
		fmt.Println(formatMessage(failure))
	}

	return nil
}

func formatMessage(failure validationFailure) string {
	moduleName, moduleType := failure.ModuleName, failure.ModuleType
	dependencyName, dependencyType := failure.DependencyName, failure.DependencyType
	return fmt.Sprintf("forbidden dependency %s(%s) → %s(%s)", moduleName, moduleType, dependencyName, dependencyType)
}
