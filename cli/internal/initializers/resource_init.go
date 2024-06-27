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
			news, err := LoadResourceFromFile(filepath.Join(resourceFolder, dirEntry.Name()))
			if err != nil {
				continue
			}
			newsSource := strings.ToLower(news[0].Source)
			newsMap[newsSource] = append(newsMap[newsSource], news...)
		}
	}
	return newsMap, nil
}
func LoadStaticResourcesFromFolder(resourceFolder string) (map[string][]string, error) {
	newsMap := make(map[string][]string)

	err := filepath.Walk(resourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileName := info.Name()

			parentDir := filepath.Dir(path)

			dirName := filepath.Base(parentDir)

			fileWithDir := filepath.Join(parentDir, fileName)
			println("Initializer path:" + dirName + " path:" + fileWithDir)

			newsMap[dirName] = append(newsMap[dirName], fileWithDir)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}

	return newsMap, nil
}
func LoadResourceFromFile(filePath string) ([]entity.News, error) {
	p, err := parser.GetFileParser(entity.PathToFile(filePath))
	if err != nil {
		return nil, err
	}
	news, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return news, nil

}
