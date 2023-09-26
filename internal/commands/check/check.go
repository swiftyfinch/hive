package check

import (
	"fmt"
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/core"
)

func Check(configPath string) error {
	config, err := config.Read(configPath)
	if err != nil {
		return err
	}
	types := core.DefaultTypes()
	if err := config.Validate(types); err != nil {
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
	rules := core.DependencyRules()
	failures, err := checkDependencies(localPods, rules, moduleTypes)
	if err != nil {
		return err
	}
	for _, failure := range failures {
		fmt.Println(formatMessage(failure))
	}

	return nil
}

func formatMessage(failure validationFailure) string {
	return fmt.Sprintf(
		"[error] %s(%s) → %s(%s)",
		failure.ModuleName,
		failure.ModuleType,
		failure.DependencyName,
		failure.DependencyType,
	)
}
