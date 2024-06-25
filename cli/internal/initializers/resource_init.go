package initializers

import (
	"fmt"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/parser"
	"os"
	"path/filepath"
	"strings"
)

func LoadResourcesFromFolder(resourceFolder string) (map[string][]entity.News, error) {
	dirEntries, err := os.ReadDir(resourceFolder)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	newsMap := make(map[string][]entity.News)
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			news, err := LoadResourceFromFile(resourceFolder, dirEntry.Name())
			if err != nil {
				continue
			}
			newsSource := strings.ToLower(news[0].Source)
			newsMap[newsSource] = append(newsMap[newsSource], news...)
		}
	}
	return newsMap, nil
}
func LoadResourceFromFile(resourceFolder, fileName string) ([]entity.News, error) {
	filePath := filepath.Join(resourceFolder, fileName)
	news, err := getResourceNews(entity.PathToFile(filePath))
	if err != nil {
		return nil, err
	}
	return news, err
}

// getResourceNews parses news from a single resource.
func getResourceNews(path entity.PathToFile) ([]entity.News, error) {
	p, err := parser.GetFileParser(path)
	if err != nil {
		return nil, err
	}
	news, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return news, nil
}
