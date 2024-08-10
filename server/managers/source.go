package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
	_ "slices"
)

// SourceManager provides API for handling news sources.
//
//go:generate mockgen -source=source.go -destination=mock_managers/mock_source.go
type SourceManager interface {
	CreateSource(name, url string) (entity.Source, error)
	GetSource(name string) (entity.Source, error)
	GetSources() ([]entity.Source, error)
	UpdateSource(name, newUrl string) error
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
	for _, source := range sources {
		if string(source.Name) == name {
			return entity.Source{}, errors.New(fmt.Sprintf("Source with name %s already exists", name))
		}
	}
	newSource := entity.Source{
		Name:       entity.SourceName(name),
		PathToFile: entity.PathToFile(url),
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
	for _, source := range sources {
		if string(source.Name) == name {
			return source, nil
		}
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
func (sourceManager sourceFolder) UpdateSource(name, newUrl string) error {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	for i, source := range sources {
		if string(source.Name) == name {
			sources[i].PathToFile = entity.PathToFile(newUrl)
			return writeToFile(sourceManager.path, sources)
		}
	}
	return fmt.Errorf("source with name %s not found", name)
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
