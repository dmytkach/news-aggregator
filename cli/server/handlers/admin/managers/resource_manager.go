package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"news-aggregator/internal/entity"
	"os"
	"regexp"
	"strings"
)

const pathToResources = "server-resources/source.json"

// GetAllSourcesNames from resource file.
func GetAllSourcesNames() ([]string, error) {
	sources, _ := readSourcesFromFile()
	resourceNames := make([]string, 0)
	for _, s := range sources {
		resourceNames = append(resourceNames, string(s.ResourceName))
	}
	if len(resourceNames) == 0 {
		return nil, errors.New("no resources found")
	}
	return resourceNames, nil
}

// GetSourcesFeeds by given name from resource file.
func GetSourcesFeeds(name string) ([]string, error) {
	sources, _ := readSourcesFromFile()
	resourceFeeds := make([]string, 0)
	for _, s := range sources {
		if strings.Contains(string(s.ResourceName), name) {
			resourceFeeds = append(resourceFeeds, string(s.PathToFile))
		}
	}
	if len(resourceFeeds) == 0 {
		return nil, errors.New("No resources found for name: " + name)
	}
	return resourceFeeds, nil
}
func CreateSource(name, url string) (entity.Resource, error) {
	sources, _ := readSourcesFromFile()
	for _, s := range sources {
		if string(s.PathToFile) == url && strings.Contains(string(s.ResourceName), name) {
			return entity.Resource{}, errors.New("resource already exists")
		}
	}
	resource := entity.Resource{
		ResourceName: entity.ResourceName(cleanFilename(name)),
		PathToFile:   entity.PathToFile(url),
	}
	sources = append(sources, resource)
	err := writeToFile(sources)
	if err != nil {
		return entity.Resource{}, err
	}
	return resource, nil
}
func RemoveSourceByName(sourceName string) error {
	sources, _ := readSourcesFromFile()
	deletedSources := make([]entity.Resource, 0)
	for _, s := range sources {
		if !strings.Contains(string(s.ResourceName), sourceName) {
			deletedSources = append(deletedSources, s)
		}
	}
	return writeToFile(deletedSources)
}
func UpdateSource(oldUrl, newUrl string) error {
	sources, err := readSourcesFromFile()
	if err != nil {
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
func writeToFile(sources []entity.Resource) error {
	jsonData, err := json.Marshal(sources)
	if err != nil {
		return err
	}

	err = os.WriteFile(pathToResources, jsonData, 0644)
	if err != nil {
		return err
	}
	return nil
}
func readSourcesFromFile() ([]entity.Resource, error) {
	file, err := os.Open(pathToResources)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
		fmt.Println("Error opening sources file:", err)
		return nil, err
	}
	defer file.Close()

	var sources []entity.Resource
	if err := json.NewDecoder(file).Decode(&sources); err != nil {
		fmt.Println("Error decoding sources file:", err)
		return nil, err
	}
	return sources, nil
}
func cleanFilename(filename string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9_]+")
	cleaned := reg.ReplaceAllString(filename, "_")
	return cleaned
}
