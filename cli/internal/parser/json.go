package parser

import (
	"encoding/json"
	"errors"
	"news-aggregator/internal/entity"
	"os"
	"time"
)

// Json - parser for JSON files.
type Json struct {
	FilePath entity.PathToFile
}

// newsResponse represents the structure of the JSON response containing news articles.
type newsResponse struct {
	Articles []newsArticle `json:"articles"`
}

// newsArticle represents a single news article parsed from JSON.
type newsArticle struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"url"`
	Date        time.Time `json:"publishedAt"`
}

// Parse - implementation of a parser for files in JSON format.
func (jsonParser *Json) Parse() ([]entity.News, error) {
	file, err := os.Open(string(jsonParser.FilePath))
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			print(err)
		}
	}(file)

	var response newsResponse
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
	if len(allNews) == 0 {
		return nil, errors.New("no news found")
	}
	return allNews, nil
}
