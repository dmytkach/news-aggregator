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

const resourceFolder = "server-news/"

var timeNow = time.Now().Format("2006-01-02")

// GetNewsFromFolder retrieves news data from a specified folder
// containing structured news resources.
func GetNewsFromFolder(folderName string) ([]entity.News, error) {
	resources, err := initializers.LoadStaticResourcesFromFolder(resourceFolder)
	if err != nil {
		log.Printf("Error loading static resources from folder: %v", err)
		return nil, err
	}
	r := resources[folderName]
	allNews := make([]entity.News, 0)
	for _, i := range r {
		news, err := GetNewsFromFile(i)
		if err != nil {
			log.Printf("Error getting news from file %s: %v", i, err)
			return nil, err
		}
		allNews = append(allNews, news...)
	}
	return allNews, nil
}

// GetNewsFromFile using parsers.
func GetNewsFromFile(filePath string) ([]entity.News, error) {
	parsers, err := parser.GetFileParser(entity.PathToFile(filePath))
	if err != nil {
		log.Printf("Error getting file parser for %s: %v", filePath, err)
		return nil, err
	}
	for _, p := range parsers {
		news, err := p.Parse()
		if err != nil {
			continue
		}
		return news, nil
	}
	log.Printf("No parsers succeeded for file %s", filePath)
	return nil, err
}

// AddNews in JSON format in the server's news folder,
// organized by source and timestamp.
func AddNews(news []entity.News, newsSource string) error {
	finalFileName := fmt.Sprintf("%s/%s.json", cleanFilename(newsSource), timeNow)
	finalFilePath := filepath.Join(resourceFolder, finalFileName)
	err := os.MkdirAll(filepath.Dir(finalFilePath), 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	jsonData, err := json.Marshal(news)
	if err != nil {
		return err
	}
	err = os.WriteFile(finalFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing news to file %s: %v", finalFilePath, err)
		return err
	}
	return nil
}
