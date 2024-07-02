package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
	"regexp"
	"strings"
)

const pathToResources = "server-resources/source.json"

// GetAllSourcesNames from resource file.
func GetAllSourcesNames() ([]string, error) {
	sources, err := readFromFile()
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return nil, err
	}
	resourceNames := make([]string, 0)
	for _, s := range sources {
		resourceNames = append(resourceNames, string(s.Name))
	}
	if len(resourceNames) == 0 {
		return nil, errors.New("no resources found")
	}
	return resourceNames, nil
}

// GetSourcesFeeds by given name from resource file.
func GetSourcesFeeds(name string) ([]string, error) {
	sources, err := readFromFile()
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return nil, err
	}
	resourceFeeds := make([]string, 0)
	for _, s := range sources {
		if strings.Contains(string(s.Name), name) {
			resourceFeeds = append(resourceFeeds, string(s.PathToFile))
		}
	}
	if len(resourceFeeds) == 0 {
		return nil, errors.New("No resources found for name: " + name)
	}
	return resourceFeeds, nil
}

// CreateSource creates a new source with the provided name and URL.
func CreateSource(name, url string) (entity.Resource, error) {
	sources, err := readFromFile()
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return entity.Resource{}, err
	}
	for _, s := range sources {
		if string(s.PathToFile) == url && strings.Contains(string(s.Name), cleanFilename(name)) {
			return entity.Resource{}, errors.New("resource already exists")
		}
	}
	resource := entity.Resource{
		Name:       entity.ResourceName(cleanFilename(name)),
		PathToFile: entity.PathToFile(url),
	}
	sources = append(sources, resource)
	err = writeToFile(sources)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return entity.Resource{}, err
	}
	log.Printf("Created new resource: %v", resource)
	return resource, nil
}

// RemoveSourceByName from the resource file.
func RemoveSourceByName(sourceName string) error {
	sources, err := readFromFile()
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	deletedSources := make([]entity.Resource, 0)
	for _, s := range sources {
		if !strings.Contains(string(s.Name), sourceName) {
			deletedSources = append(deletedSources, s)
		}
	}
	err = writeToFile(deletedSources)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return err
	}
	log.Printf("Removed source with name: %s", sourceName)
	return nil
}

// UpdateSource identified by its old URL.
func UpdateSource(oldUrl, newUrl string) error {
	sources, err := readFromFile()
	if err != nil {
		log.Printf("Error reading from file: %v", err)
		return err
	}
	for _, s := range sources {
		if string(s.PathToFile) == newUrl {
			return errors.New("resource already exists")
		}
	}
	for i, s := range sources {
		if string(s.PathToFile) == oldUrl {
			sources[i].PathToFile = entity.PathToFile(newUrl)
			return writeToFile(sources)
		}
	}
	return fmt.Errorf("source with URL %s not found", oldUrl)
}

// writeToFile sources in JSON format.
func writeToFile(sources []entity.Resource) error {
	jsonData, err := json.Marshal(sources)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		return err
	}

	err = os.WriteFile(pathToResources, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing to file: %v", err)
		return err
	}
	return nil
}

// readFromFile resources file.
func readFromFile() ([]entity.Resource, error) {
	file, err := os.Open(pathToResources)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Resource file does not exist: %v", err)
			return nil, err
		}
		log.Printf("Error opening sources file: %v", err)
		return nil, err
	}
	defer file.Close()

	var sources []entity.Resource
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		log.Printf("Error decoding sources file: %v", err)
		return nil, err
	}
	return sources, nil
}
func cleanFilename(filename string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	cleaned := reg.ReplaceAllString(filename, "_")
	return cleaned
}
