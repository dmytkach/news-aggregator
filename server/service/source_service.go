package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

// AddSource a new news source whose url was transmitted to the response.
// The addition occurs by the name and url of the source.
func AddSource(url string) (entity.Source, error) {
	feed, err := fetchFeed(entity.PathToFile(url))
	if err != nil {
		log.Print("error loading feed")
		return entity.Source{}, err
	}
	return managers.CreateSource(string(feed.Name), url)
}

// GetSources names of registered news sources.
func GetSources() ([]entity.Source, error) {
	return managers.GetSources()
}

// GetSource url by of given news name.
func GetSource(name string) (entity.Source, error) {
	return managers.GetSource(name)
}

// UpdateSource the URL of an existing news source.
func UpdateSource(oldUrl, newUrl string) error {
	return managers.UpdateSource(oldUrl, newUrl)
}

// RemoveSource a news source by its name.
func RemoveSource(newsName string) error {
	return managers.RemoveSourceByName(newsName)
}
