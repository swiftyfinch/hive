package cocoapods

import (
	"fmt"
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
	remotePods map[string]Pod,
	localPods map[string]Pod,
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

	localPods = map[string]Pod{}
	remotePods = map[string]Pod{}
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

type Pod struct {
	Name         string
	Dependencies []string
}

func parsePods(anyPods []interface{}) (map[string]Pod, error) {
	pods := map[string]Pod{}
	for _, any := range anyPods {
		switch pod := any.(type) {
		case string:
			name, _, err := parsePod(pod)
			if err != nil {
				return nil, err
			}
			pods[name] = Pod{name, []string{}}
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
				pods[name] = Pod{name, dependencies}
			}
		}
	}
	return pods, nil
}

func parsePod(line string) (string, string, error) {
	re := regexp.MustCompile(`([^\s]*)(?:\s\((.*)\))?`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return "", "", fmt.Errorf("incorrect pod name and version format")
	}
	return matches[1], matches[2], nil
}
