package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
)

// SourceManager provides API for handling news sources.
type SourceManager interface {
	CreateSource(name, url string) (entity.Source, error)
	GetSource(name string) (entity.Source, error)
	GetSources() ([]entity.Source, error)
	UpdateSource(oldUrl, newUrl string) error
	RemoveSourceByName(sourceName string) error
}

// sourceFolder implements SourceManager using a folder-based storage for sources.
type sourceFolder struct {
	path string
}

// CreateSourceFolder in the specified directory.
func CreateSourceFolder(pathToSources string) SourceManager {
	return sourceFolder{pathToSources}
}

// CreateSource creates a new source with the provided name and URL.
func (sourceManager sourceFolder) CreateSource(name, url string) (entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return entity.Source{}, err
	}

	existingSource, found := findSourceByName(sources, name)
	if found {
		err = updateSourceWithPath(&existingSource, url, sources, sourceManager.path)
		if err != nil {
			return entity.Source{}, err
		}
		log.Printf("Updated resource: %v", existingSource)
		return existingSource, nil
	}

	newSource := entity.Source{
		Name:        entity.SourceName(name),
		PathsToFile: []entity.PathToFile{entity.PathToFile(url)},
	}
	sources = append(sources, newSource)
	err = writeToFile(sourceManager.path, sources)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return entity.Source{}, err
	}
	log.Printf("Created new resource: %v", newSource)
	return newSource, nil
}

// GetSource by given name from resource file.
func (sourceManager sourceFolder) GetSource(name string) (entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return entity.Source{}, err
	}
	existingSource, found := findSourceByName(sources, name)
	if found {
		return existingSource, nil
	}
	return entity.Source{}, errors.New("no resources found for name: " + name)
}

// GetSources from source file.
func (sourceManager sourceFolder) GetSources() ([]entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return nil, err
	}
	return sources, nil
}

// UpdateSource identified by its old URL.
func (sourceManager sourceFolder) UpdateSource(oldUrl, newUrl string) error {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	for i, s := range sources {
		for j, path := range s.PathsToFile {
			if string(path) == oldUrl {
				if isPathExist(s.PathsToFile, newUrl) {
					return errors.New("resource already exists")
				}
				sources[i].PathsToFile[j] = entity.PathToFile(newUrl)
				return writeToFile(sourceManager.path, sources)
			}
		}
	}
	return fmt.Errorf("source with URL %s not found", oldUrl)
}

// RemoveSourceByName from the resource file.
func (sourceManager sourceFolder) RemoveSourceByName(sourceName string) error {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	deletedSources := make([]entity.Source, 0)
	for _, s := range sources {
		if string(s.Name) != sourceName {
			deletedSources = append(deletedSources, s)
		}
	}
	err = writeToFile(sourceManager.path, deletedSources)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return err
	}
	log.Printf("Removed source with name: %s", sourceName)
	return nil
}

func findSourceByName(sources []entity.Source, name string) (entity.Source, bool) {
	for _, s := range sources {
		if string(s.Name) == name {
			return s, true
		}
	}
	return entity.Source{}, false
}

func updateSourceWithPath(source *entity.Source, url string, sources []entity.Source, path string) error {
	if isPathExist(source.PathsToFile, url) {
		return errors.New("resource already exists")
	}
	source.PathsToFile = append(source.PathsToFile, entity.PathToFile(url))
	return writeToFile(path, sources)
}

func isPathExist(paths []entity.PathToFile, url string) bool {
	for _, path := range paths {
		if string(path) == url {
			return true
		}
	}
	return false
}

// writeToFile sources in JSON format.
func writeToFile(path string, sources []entity.Source) error {
	jsonData, err := json.MarshalIndent(sources, "", "  ")
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return err
	}

	err = os.WriteFile(path, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return err
	}
	return nil
}

// readFromFile resources file.
func readFromFile(path string) ([]entity.Source, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Source file does not exist, creating a new one: %v", path)
			newFile, err := os.Create(path)
			if err != nil {
				log.Printf("Error creating new source file: %v", err)
				return nil, err
			}
			defer func(newFile *os.File) {
				err := newFile.Close()
				if err != nil {
					log.Printf("failed to close file: %v", err)
				}
			}(newFile)
			var emptySources []entity.Source
			if err := writeToFile(path, emptySources); err != nil {
				log.Printf("Error initializing new source file: %v", err)
				return nil, err
			}
			return emptySources, nil
		}
		log.Printf("Error opening sources file: %v", err)
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}(file)

	var sources []entity.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		log.Printf("Error decoding sources file: %v", err)
		return nil, err
	}
	return sources, nil
}
