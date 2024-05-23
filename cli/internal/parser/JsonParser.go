package parser

import (
	"NewsAggregator/cli/internal/entity"
	"encoding/json"
	"os"
	"time"
)

// JsonParser - parser for JSON files.
type JsonParser struct{}

// NewsResponse represents the structure of the JSON response containing news articles.
type NewsResponse struct {
	Articles []NewsArticle `json:"articles"`
}

// NewsArticle represents a single news article parsed from JSON.
type NewsArticle struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"url"`
	Date        time.Time `json:"publishedAt"`
}

// Parse - implementation of a parser for files in JSON format.
func (jsonParser *JsonParser) Parse(FileName string) ([]entity.News, error) {
	file, err := os.Open(FileName)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			print(err)
		}
	}(file)

	var response NewsResponse
	err = json.NewDecoder(file).Decode(&response)
	if err != nil {
		return nil, err
	}

	var allNews []entity.News

	for _, article := range response.Articles {
		news := entity.News{
			Title:       entity.Title(article.Title),
			Description: entity.Description(article.Description),
			Link:        entity.Link(article.Link),
			Date:        article.Date,
		}
		allNews = append(allNews, news)
	}
	return allNews, nil
}
