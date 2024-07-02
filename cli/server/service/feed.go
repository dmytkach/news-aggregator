package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

// Add a new news source whose url was transmitted to the response.
// The addition occurs by the name and url of the source.
func Add(url string) (entity.Resource, error) {
	news, err := fetchNewsFromResponse(entity.PathToFile(url))
	if err != nil {
		log.Print("error loading feed")
		return entity.Resource{}, err
	}
	name := news[0].Source
	return managers.CreateSource(name, url)
}

// GetAll names of registered news sources.
func GetAll() ([]string, error) {
	return managers.GetAllSourcesNames()
}

// Get all url by of given news name.
func Get(newsName string) ([]string, error) {
	return managers.GetSourcesFeeds(newsName)
}

// Update the URL of an existing news source.
func Update(oldUrl, newUrl string) error {
	return managers.UpdateSource(oldUrl, newUrl)
}

// Remove a news source by its name.
func Remove(newsName string) error {
	return managers.RemoveSourceByName(newsName)
}
