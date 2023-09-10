package glob

import (
	"os"
	"path/filepath"
)

func FindPathsRecursively(dir string, filename string) ([]string, error) {
	paths := []string{}
	err := filepath.Walk(dir, func(path string, _ os.FileInfo, err error) error {
		if filepath.Base(path) == filename {
			paths = append(paths, path)
		}
		return nil
	})
	return paths, err
}
