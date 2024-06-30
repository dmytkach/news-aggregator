package managers

import (
	"encoding/json"
	"fmt"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/initializers"
	"news-aggregator/internal/parser"
	"os"
	"path/filepath"
	"time"
)

const resourceFolder = "server-news/"

var timeNow = time.Now().Format("2006-01-02")

func GetNewsFromFolder(folderName string) ([]entity.News, error) {
	resources, err := initializers.LoadStaticResourcesFromFolder(resourceFolder)
	if err != nil {
		return nil, err
	}
	r := resources[folderName]
	allNews := make([]entity.News, 0)
	for i := range r {
		news, err := GetNewsFromFile(r[i])
		if err != nil {
			return nil, err
		}
		allNews = append(allNews, news...)
	}
	return allNews, nil
}
func GetNewsFromFile(filePath string) ([]entity.News, error) {
	parsers, err := parser.GetFileParser(entity.PathToFile(filePath))
	if err != nil {
		return nil, err
	}
	for _, p := range parsers {
		news, err := p.Parse()
		if err != nil {
			continue
		}
		return news, nil
	}
	return nil, err
}
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
		return err
	}
	return nil
}
