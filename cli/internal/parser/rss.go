package parser

import (
	"NewsAggregator/cli/internal/entity"
	"github.com/mmcdole/gofeed"
	"os"
)

// Rss - parser for RSS files.
type Rss struct {
	FilePath entity.PathToFile
}

// Parse - implementation of a parser for files in RSS format.
func (rssParser *Rss) Parse() ([]entity.News, error) {
	file, err := os.Open(string(rssParser.FilePath))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		return nil, err
	}

	var newsItems []entity.News
	for _, item := range feed.Items {
		newsItems = append(newsItems, entity.News{
			Title:       entity.Title(item.Title),
			Description: entity.Description(item.Description),
			Link:        entity.Link(item.Link),
			Date:        *item.PublishedParsed,
		})
	}

	return newsItems, nil
}
