package parser

import (
	"NewsAggregator/cli/internal/entity"
	"github.com/PuerkitoBio/goquery"
	"os"
	"time"
)

// HtmlParser - parser for HTML files.
type HtmlParser struct {
}

// Parse - implementation of a parser for files in HTML format.
func (htmlParser *HtmlParser) Parse(FileName string) ([]entity.News, error) {
	file, err := os.Open(FileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, err
	}

	var newsItems []entity.News
	doc.Find(".news-item").Each(func(i int, s *goquery.Selection) {
		title := s.Find(".news-title").Text()
		description := s.Find(".news-description").Text()
		link, _ := s.Find(".news-link").Attr("href")
		dateStr := s.Find(".news-date").Text()
		date, _ := time.Parse("2006-01-02", dateStr)

		newsItems = append(newsItems, entity.News{
			Title:       entity.Title(title),
			Description: entity.Description(description),
			Link:        entity.Link(link),
			Date:        date,
		})
	})

	return newsItems, nil
}
