package initializers

import (
	"fmt"
	"os"
	"path/filepath"
)

// LoadSources walks through a directory and loads static sources.
// It returns a map where keys are resource names
// and values are slices of file paths for each resource.
func LoadSources(sourceFolder string) (map[string][]string, error) {
	newsMap := make(map[string][]string)
	err := filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			resourceName := filepath.Base(filepath.Dir(path))

			pathToFile := filepath.Join(filepath.Dir(path), info.Name())

			newsMap[resourceName] = append(newsMap[resourceName], pathToFile)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return newsMap, nil
}
