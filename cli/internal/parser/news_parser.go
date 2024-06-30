package parser

import (
	"encoding/json"
	"fmt"
	"news-aggregator/internal/entity"
	"os"
)

type NewsParser struct {
	FilePath entity.PathToFile
}

func (newsParser *NewsParser) Parse() ([]entity.News, error) {
	file, err := os.Open(string(newsParser.FilePath))
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
		}
	}(file)
	var articles []entity.News
	if err := json.NewDecoder(file).Decode(&articles); err != nil {
		return nil, err
	}
	return articles, nil
}

// CanParseFileType checks if the file extension is .json
func (newsParser *NewsParser) CanParseFileType(ext string) bool {
	return ext == ".json"
}
