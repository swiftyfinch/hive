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
	var severity string
	if failure.IsWarning {
		severity = "warning"
	} else {
		severity = "error"
	}
	return fmt.Sprintf(
		"[%s] %s(%s) → %s(%s)",
		severity,
		failure.ModuleName,
		failure.ModuleType,
		failure.DependencyName,
		failure.DependencyType,
	)
}
