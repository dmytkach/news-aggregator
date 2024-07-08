package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/parser"
	"os"
	"path/filepath"
	"time"
)

var NewsFolder = "server-news/"

var timeNow = time.Now().Format("2006-01-02")

// GetNewsFromFolder retrieves news data from a specified folder
// containing structured news resources.
func GetNewsFromFolder(folderName string) ([]entity.News, error) {
	resources, err := initializers.LoadSources(NewsFolder)
	if err != nil {
		log.Printf("Error loading static resources from folder: %v", err)
		return nil, err
	}
	r := resources[folderName]
	allNews := make([]entity.News, 0)
	for _, path := range r {
		file, err := os.Open(path)
		if err != nil {
			log.Printf("Failed to open file %s: %v", path, err)
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
			log.Printf("Error decoding file %s: %v", path, err)
			return nil, err
		}
		allNews = append(allNews, articles...)
	}
	return allNews, nil
}

// GetNewsFromFile using parsers.
func GetNewsFromFile(filePath string) (entity.Feed, error) {
	p, err := parser.GetFileParser(entity.PathToFile(filePath))
	if err != nil {
		log.Printf("Error getting file parser for %s: %v", filePath, err)
		return entity.Feed{}, err
	}
	f, err := p.Parse()
	if err != nil {
		log.Printf("Error parsing file %s: %v", filePath, err)
		return entity.Feed{}, err
	}
	return f, err
}

// AddNews in JSON format in the server's news folder,
// organized by source and timestamp.
func AddNews(newsToAdd []entity.News, newsSource string) error {
	finalFileName := fmt.Sprintf("%s/%s.json", newsSource, timeNow)
	finalFilePath := filepath.Join(NewsFolder, finalFileName)
	err := os.MkdirAll(filepath.Dir(finalFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	currentNews, err := loadNewsFromFile(finalFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			return fmt.Errorf("failed to load current news: %w", err)
		}
		currentNews = []entity.News{}
	}
	currentNews = append(currentNews, newsToAdd...)
	jsonData, err := json.Marshal(currentNews)
	if err != nil {
		return fmt.Errorf("failed to marshal news to JSON: %w", err)
	}
	err = os.WriteFile(finalFilePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write news to file: %w", err)
	}
	return nil
}

func loadNewsFromFile(filePath string) ([]entity.News, error) {
	var news []entity.News
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonData, &news)
	if err != nil {
		return nil, err
	}
	return news, nil
}
