package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

type NewsFeed struct {
}

func (newsFeed NewsFeed) Add(url string) (entity.Resource, error) {
	news, err := FetchNewsFromResponse(entity.PathToFile(url))
	if err != nil {
		log.Print("error loading feed")
		return entity.Resource{}, err
	}
	name := news[0].Source
	return managers.CreateSource(name, url)
}
func (newsFeed NewsFeed) GetAll() ([]string, error) {
	return managers.GetAllSourcesNames()
}
func (newsFeed NewsFeed) Get(newsName string) ([]string, error) {
	return managers.GetSourcesFeeds(newsName)
}
func (newsFeed NewsFeed) Update(oldUrl, newUrl string) error {
	return managers.UpdateSource(oldUrl, newUrl)
}
func (newsFeed NewsFeed) Remove(newsName string) error {
	return managers.RemoveSourceByName(newsName)
}
