package parser

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"news-aggregator/internal/entity"
	"os"
	"strings"
)

// Rss - parser for RSS files.
type Rss struct {
	FilePath entity.PathToFile
}

func (rssParser *Rss) CanParseFileType(ext string) bool {
	return ext == ".xml" || ext == ".rss"
}

// Parse - implementation of a parser for files in RSS format.
func (rssParser *Rss) Parse() ([]entity.News, error) {
	file, err := os.Open(string(rssParser.FilePath))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
			return
		}
	}(file)

	fp := gofeed.NewParser()
	feed, err := fp.Parse(file)
	if err != nil {
		return nil, err
	}

	var allNews []entity.News
	for _, item := range feed.Items {
		allNews = append(allNews, entity.News{
			Title:       entity.Title(item.Title),
			Description: entity.Description(item.Description),
			Link:        entity.Link(item.Link),
			Date:        *item.PublishedParsed,
			Source:      strings.ToLower(feed.Title),
		})
	}
	if len(allNews) == 0 {
		return nil, errors.New("no news found")
	}
	return allNews, nil
}
