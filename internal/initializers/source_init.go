package initializers

import (
	"fmt"
	"log"
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
			log.Print(resourceName)
			pathToFile := filepath.Join(filepath.Dir(path), info.Name())
			log.Print(pathToFile)
			newsMap[resourceName] = append(newsMap[resourceName], pathToFile)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	return newsMap, nil
}

// AnalyzeDirectory analyzes the contents of a given directory.
// It returns a slice of strings, each representing a full path to a file or directory.
func AnalyzeDirectory(directory string) ([]string, error) {
	var entries []string
	file, err := os.Open(directory)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(file)

	files, err := file.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, f := range files {
		fullPath := filepath.Join(directory, f.Name())
		entries = append(entries, fullPath)
	}

	return entries, nil
}
