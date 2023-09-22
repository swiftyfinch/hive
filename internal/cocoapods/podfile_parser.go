package cocoapods

import (
	"fmt"
	"main/internal/modules"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type podfile struct {
	Pods            []interface{}          `yaml:"PODS"`
	SpecRepos       map[string][]string    `yaml:"SPEC REPOS"`
	CheckoutOptions map[string]interface{} `yaml:"CHECKOUT OPTIONS"`
}

func ParsePodfile(path string) (
	remotePods map[string]modules.Module,
	localPods map[string]modules.Module,
	err error,
) {
	buffer, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	podfile := &podfile{}
	if err = yaml.Unmarshal(buffer, podfile); err != nil {
		return nil, nil, err
	}

	pods, err := parsePods(podfile.Pods)
	if err != nil {
		return nil, nil, err
	}

	type void struct{}
	remotePodNames := map[string]void{}
	for _, pods := range podfile.SpecRepos {
		for _, name := range pods {
			remotePodNames[name] = void{}
		}
	}
	for name := range podfile.CheckoutOptions {
		remotePodNames[name] = void{}
	}

	localPods = map[string]modules.Module{}
	remotePods = map[string]modules.Module{}
	for _, pod := range pods {
		podName := strings.Split(pod.Name, "/")[0]
		if _, ok := remotePodNames[podName]; ok {
			remotePods[pod.Name] = pod
		} else {
			localPods[pod.Name] = pod
		}
	}

	return remotePods, localPods, nil
}

func parsePods(anyPods []interface{}) (map[string]modules.Module, error) {
	pods := map[string]modules.Module{}
	for _, any := range anyPods {
		switch pod := any.(type) {
		case string:
			name, _, err := parsePod(pod)
			if err != nil {
				return nil, err
			}
			pods[name] = modules.Module{Name: name, Dependencies: []string{}}
		case map[interface{}]interface{}:
			for anyName, anyDependencies := range pod {
				line, ok := anyName.(string)
				if !ok {
					return nil, fmt.Errorf("incorrect type of pod name")
				}
				name, _, err := parsePod(line)
				if err != nil {
					return nil, err
				}

				dependencies, err := parseDependencies(anyDependencies)
				if err != nil {
					return nil, err
				}
				pods[name] = modules.Module{Name: name, Dependencies: dependencies}
			}
		case map[string]interface{}:
			for line, anyDependencies := range pod {
				name, _, err := parsePod(line)
				if err != nil {
					return nil, err
				}

				dependencies, err := parseDependencies(anyDependencies)
				if err != nil {
					return nil, err
				}
				pods[name] = modules.Module{Name: name, Dependencies: dependencies}
			}
		default:
			return nil, fmt.Errorf("unkown type of pods")
		}
	}
	return pods, nil
}

func parseDependencies(anyDependencies interface{}) ([]string, error) {
	anyDependenciesSlice, ok := anyDependencies.([]interface{})
	if !ok {
		return nil, fmt.Errorf("incorrect type of pod dependencies")
	}

	dependencies := []string{}
	for _, anyDependency := range anyDependenciesSlice {
		dependency, ok := anyDependency.(string)
		if !ok {
			return nil, fmt.Errorf("incorrect type of pod dependency")
		}

		name, _, err := parsePod(dependency)
		if err != nil {
			return nil, err
		}
		dependencies = append(dependencies, name)
	}
	return dependencies, nil
}

func parsePod(line string) (string, string, error) {
	re := regexp.MustCompile(`([^\s]*)(?:\s\((.*)\))?`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("incorrect pod name and version format")
	}
	return matches[1], matches[2], nil
}
