package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

func Add(url string) (entity.Resource, error) {
	news, err := fetchNewsFromResponse(entity.PathToFile(url))
	if err != nil {
		log.Print("error loading feed")
		return entity.Resource{}, err
	}
	name := news[0].Source
	return managers.CreateSource(name, url)
}
func GetAll() ([]string, error) {
	return managers.GetAllSourcesNames()
}
func Get(newsName string) ([]string, error) {
	return managers.GetSourcesFeeds(newsName)
}
func Update(oldUrl, newUrl string) error {
	return managers.UpdateSource(oldUrl, newUrl)
}
func Remove(newsName string) error {
	return managers.RemoveSourceByName(newsName)
}
