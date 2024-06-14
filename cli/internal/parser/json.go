package parser

import (
	"encoding/json"
	"errors"
	"fmt"
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

func (jsonParser *Json) CanParseFileType(ext string) bool {
	return ext == ".json"
}

// Parse - implementation of a parser for files in JSON format.
func (jsonParser *Json) Parse() ([]entity.News, error) {
	file, err := os.Open(string(jsonParser.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
			return
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
