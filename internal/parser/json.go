package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"news-aggregator/internal/entity"
	"os"
	"regexp"
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

// Parse - implementation of a parser for files in JSON format.
func (jsonParser *Json) Parse() (entity.Feed, error) {
	file, err := os.Open(string(jsonParser.FilePath))
	if err != nil {
		return entity.Feed{}, fmt.Errorf("failed to open file: %w", err)
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
		return entity.Feed{}, err
	}

	var allNews []entity.News
	var title string
	for _, article := range response.Articles {
		news := entity.News{
			Title:       entity.Title(article.Title),
			Description: entity.Description(article.Description),
			Link:        entity.Link(article.Link),
			Date:        article.Date,
			Source:      cleanSourceName(article.Source.Name),
		}
		title = cleanSourceName(article.Source.Name)
		allNews = append(allNews, news)
	}
	if len(allNews) == 0 {
		return entity.Feed{}, errors.New("no news found")
	}
	return entity.Feed{Name: entity.SourceName(title), News: allNews}, nil
}

// CanParseFileType checks if the file extension is .json
func (jsonParser *Json) CanParseFileType(ext string) bool {
	return ext == ".json"
}
func cleanSourceName(filename string) string {
	reg := regexp.MustCompile(`[^\p{L}\p{N}_]+`)
	cleaned := reg.ReplaceAllString(filename, "_")
	return strings.ToLower(cleaned)
}
