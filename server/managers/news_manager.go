package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/initializers"
	"os"
	"path/filepath"
	"time"
)

var timeNow = time.Now().Format("2006-01-02")

type NewsManager interface {
	AddNews(newsToAdd []entity.News, newsSource string) error
	GetNewsFromFolder(folderName string) ([]entity.News, error)
	GetNewsSourceFilePath(sourceName []string) (map[string][]string, error)
}

type newsFolderManager struct {
	path string
}

func CreateNewsFolderManager(pathToNews string) NewsManager {
	return newsFolderManager{pathToNews}
}

// AddNews in JSON format in the server's news folder,
// organized by source and timestamp.
func (repo newsFolderManager) AddNews(newsToAdd []entity.News, newsSource string) error {
	finalFileName := fmt.Sprintf("%s/%s.json", newsSource, timeNow)
	finalFilePath := filepath.Join(repo.path, finalFileName)
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

// GetNewsFromFolder retrieves news data from a specified folder
// containing structured news resources.
func (repo newsFolderManager) GetNewsFromFolder(folderName string) ([]entity.News, error) {
	sourcePath := filepath.Join(repo.path, folderName)
	resources, err := getNewsSources(sourcePath)
	if err != nil {
		return nil, err
	}
	allNews := make([]entity.News, 0)
	for _, path := range resources {
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

func (repo newsFolderManager) GetNewsSourceFilePath(sourceName []string) (map[string][]string, error) {
	resources := make(map[string][]string)
	for _, source := range sourceName {
		sourcePath := filepath.Join(repo.path, source)
		paths, err := getNewsSources(sourcePath)
		if err != nil {
			log.Print("Not found news source")
			return nil, err
		}
		resources[source] = paths
	}
	return resources, nil
}

func getNewsSources(sourceName string) ([]string, error) {
	s, err := initializers.AnalyzeDirectory(sourceName)
	if err != nil {
		log.Printf("Failed to analyze source %s: %v", sourceName, err)
		return nil, err
	}
	return s, nil
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
