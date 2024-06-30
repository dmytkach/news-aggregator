package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"news-aggregator/internal/entity"
	"os"
	"strings"
	"time"
)

// Json represents a JSON parser for news articles.
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
	Source      struct {
		Name string `json:"name"`
	} `json:"source"`
}

// Parse implements a parser for JSON files, attempting to parse into different structures.
func (jsonParser *Json) Parse() ([]entity.News, error) {
	file, err := os.Open(string(jsonParser.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
		}
	}(file)

	var allNews []entity.News

	if news, err := jsonParser.decodeNewsResponse(file); err == nil {
		allNews = append(allNews, news...)
	} else {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	if len(allNews) == 0 {
		return nil, errors.New("no news found")
	}

	return allNews, nil
}

// decodeNewsResponse attempts to decode the file into the newsResponse structure.
func (jsonParser *Json) decodeNewsResponse(file *os.File) ([]entity.News, error) {
	var response newsResponse
	if err := json.NewDecoder(file).Decode(&response); err != nil {
		return nil, err
	}

	var allNews []entity.News
	for _, article := range response.Articles {
		news := entity.News{
			Title:       entity.Title(article.Title),
			Description: entity.Description(article.Description),
			Link:        entity.Link(article.Link),
			Date:        article.Date,
			Source:      strings.ToLower(article.Source.Name),
		}
		allNews = append(allNews, news)
	}

	return allNews, nil
}

// CanParseFileType checks if the file extension is .json
func (jsonParser *Json) CanParseFileType(ext string) bool {
	return ext == ".json"
}
