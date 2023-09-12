package cocoapods

import (
	"fmt"
	glob "hive/packages/utils"
)

func ReadPods() (
	remotePods map[string]Pod,
	localPods map[string]Pod,
	err error,
) {
	paths, err := glob.FindPathsRecursively(".", "Podfile.lock")
	if err != nil {
		return nil, nil, err
	}
	if len(paths) == 0 {
		return nil, nil, fmt.Errorf("couldn't find any Podfile.lock")
	} else if len(paths) > 1 {
		return nil, nil, fmt.Errorf("found several Podfile.lock files")
	}
	return ParsePodfile(paths[0])
}
