package config

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

type Ignore struct {
	ModuleRegexp       regexp.Regexp
	DependenciesRegexp regexp.Regexp
}

func ReadIgnore(path string) ([]Ignore, error) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	ignoreYML := &[]map[string]string{}
	if err = yaml.Unmarshal(buffer, ignoreYML); err != nil {
		return nil, err
	}

	ignores := []Ignore{}
	for _, ignore := range *ignoreYML {
		for key, value := range ignore {
			moduleRegexp, err := regexp.Compile(key)
			if err != nil {
				continue
			}

			dependenciesRegexp, err := regexp.Compile(value)
			if err != nil {
				continue
			}

			ignores = append(
				ignores,
				Ignore{ModuleRegexp: *moduleRegexp, DependenciesRegexp: *dependenciesRegexp},
			)
		}
	}
	return ignores, nil
}
