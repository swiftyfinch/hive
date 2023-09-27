package check

import (
	"fmt"
	"main/internal/cocoapods"
	"main/internal/config"
	"main/internal/core"
	"os"
)

func Check(workingDirectory string) error {
	configPath := workingDirectory + "/" + core.Modules_File_Name
	config, err := config.Read(configPath)
	if err != nil {
		return err
	}

	types := core.DefaultTypes()
	if err := config.Validate(types); err != nil {
		return err
	}

	ignore, err := readIgnore(workingDirectory + "/ignore.yml")
	if err != nil {
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

	failures, err := checkDependencies(localPods, rules, moduleTypes, ignore)
	if err != nil {
		return err
	}
	for _, failure := range failures {
		fmt.Println(formatMessage(failure))
	}

	return nil
}

func readIgnore(path string) ([]config.Ignore, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return []config.Ignore{}, nil
	}
	return config.ReadIgnore(path)
}

func formatMessage(failure validationFailure) string {
	return fmt.Sprintf(
		"⛔️ [%s: %s] %s → %s",
		failure.ModuleType,
		failure.DependencyType,
		failure.ModuleName,
		failure.DependencyName,
	)
}
