package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
)

type SourceManager interface {
	CreateSource(name, url string) (entity.Source, error)
	GetSource(name string) (entity.Source, error)
	GetSources() ([]entity.Source, error)
	UpdateSource(oldUrl, newUrl string) error
	RemoveSourceByName(sourceName string) error
}

type sourceFolderManager struct {
	path string
}

func CreateSourceFolderManager(pathToSources string) SourceManager {
	return sourceFolderManager{pathToSources}
}

// GetSources from source file.
func (sourceManager sourceFolderManager) GetSources() ([]entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return nil, err
	}
	return sources, nil
}

// GetSource by given name from resource file.
func (sourceManager sourceFolderManager) GetSource(name string) (entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return entity.Source{}, err
	}
	for _, s := range sources {
		if string(s.Name) == name {
			return s, nil
		}
	}
	return entity.Source{}, errors.New("No resources found for name: " + name)
}

// CreateSource creates a new source with the provided name and URL.
func (sourceManager sourceFolderManager) CreateSource(name, url string) (entity.Source, error) {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return entity.Source{}, err
	}
	for i, s := range sources {
		if string(s.Name) == name {
			for _, path := range s.PathsToFile {
				if string(path) == url {
					return entity.Source{}, errors.New("resource already exists")
				}
			}
			sources[i].PathsToFile = append(sources[i].PathsToFile, entity.PathToFile(url))
			err = writeToFile(sourceManager.path, sources)
			if err != nil {
				log.Printf("Error writing to file: %v", err)
				return entity.Source{}, err
			}
			log.Printf("Updated resource: %v", sources[i])
			return sources[i], nil
		}
	}
	resource := entity.Source{
		Name:        entity.SourceName(name),
		PathsToFile: []entity.PathToFile{entity.PathToFile(url)},
	}
	sources = append(sources, resource)
	err = writeToFile(sourceManager.path, sources)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return entity.Source{}, err
	}
	log.Printf("Created new resource: %v", resource)
	return resource, nil
}

// RemoveSourceByName from the resource file.
func (sourceManager sourceFolderManager) RemoveSourceByName(sourceName string) error {
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
	log.Printf("Removed source with name: %sourceManager", sourceName)
	return nil
}

// UpdateSource identified by its old URL.
func (sourceManager sourceFolderManager) UpdateSource(oldUrl, newUrl string) error {
	sources, err := readFromFile(sourceManager.path)
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	for i, s := range sources {
		for j, path := range s.PathsToFile {
			if string(path) == oldUrl {
				for _, p := range s.PathsToFile {
					if string(p) == newUrl {
						return errors.New("resource already exists")
					}
				}
				sources[i].PathsToFile[j] = entity.PathToFile(newUrl)
				return writeToFile(sourceManager.path, sources)
			}
		}
	}
	return fmt.Errorf("source with URL %s not found", oldUrl)
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
			defer newFile.Close()

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
	defer file.Close()

	var sources []entity.Source
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		log.Printf("Error decoding sources file: %v", err)
		return nil, err
	}
	return sources, nil
}
