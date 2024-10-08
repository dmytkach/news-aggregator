package managers

import (
	"encoding/json"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"os"
	"path/filepath"
	"time"
)

var timeNow = time.Now().Format("2006-01-02")

// NewsManager provides API for handling news data.
//
//go:generate mockgen -source=news.go -destination=mock_managers/mock_news.go
type NewsManager interface {
	AddNews(newsToAdd []entity.News, newsSource string) error
	GetNewsFromFolder(folderName string) ([]entity.News, error)
	GetNewsSourceFilePath(sourceName []string) (map[string][]string, error)
}

// newsFolder implements the NewsManager for managing news data stored in folders.
type newsFolder struct {
	path string
}

// CreateNewsFolder with the given path.
func CreateNewsFolder(pathToNews string) NewsManager {
	return newsFolder{pathToNews}
}

// AddNews in JSON format in the server's news folder,
// organized by source and timestamp.
func (folder newsFolder) AddNews(newsToAdd []entity.News, newsSource string) error {
	finalFileName := fmt.Sprintf("%s/%s.json", newsSource, timeNow)
	finalFilePath := filepath.Join(folder.path, finalFileName)
	err := os.MkdirAll(filepath.Dir(finalFilePath), 0755)
	log.Printf("Final file path: %s", finalFilePath)
	if err != nil {
		log.Printf("Error creating directory: %v", err)
		return fmt.Errorf("failed to create directory: %w", err)
	}
	currentNews, err := loadNewsFromFile(finalFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Error loading current news: %v", err)
			return fmt.Errorf("failed to load current news: %w", err)
		}
		log.Printf("File does not exist, сreate a new file ...")
		currentNews = []entity.News{}
	}
	currentNews = append(currentNews, newsToAdd...)
	jsonData, err := json.Marshal(currentNews)
	if err != nil {
		log.Printf("Error marshalling news to JSON: %v", err)
		return fmt.Errorf("failed to marshal news to JSON: %w", err)
	}
	err = os.WriteFile(finalFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("Error writing news to file: %v", err)
		return fmt.Errorf("failed to write news to file: %w", err)
	}
	log.Printf("Successfully added %d news items to %s", len(newsToAdd), finalFilePath)
	return nil
}

// GetNewsFromFolder retrieves news data from a specified folder
// containing structured news resources.
func (folder newsFolder) GetNewsFromFolder(folderName string) ([]entity.News, error) {
	sourcePath := filepath.Join(folder.path, folderName)
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

// GetNewsSourceFilePath provides the file paths associated with each news source,
// helping in locating where news data for each source is stored.
func (folder newsFolder) GetNewsSourceFilePath(sourceName []string) (map[string][]string, error) {
	resources := make(map[string][]string)
	for _, s := range sourceName {
		sourcePath := filepath.Join(folder.path, s)
		paths, err := getNewsSources(sourcePath)
		if err != nil {
			log.Print("Not found news")
			return nil, err
		}
		if len(paths) != 0 {
			resources[s] = paths
			log.Printf("Found news paths for %s", s)
		}
	}
	return resources, nil
}

// getNewsSources analyzes the contents of a given directory.
// It returns a slice of a full paths to a file .
func getNewsSources(sourceName string) ([]string, error) {
	_, err := os.Stat(sourceName)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(sourceName, os.ModePerm); err != nil {
			return nil, fmt.Errorf("failed to create directory: %w", err)
		}
		log.Printf("Directory created: %s", sourceName)
	}

	dir, err := os.Open(sourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open directory: %w", err)
	}
	defer func() {
		if err := dir.Close(); err != nil {
			log.Printf("failed to close directory: %v", err)
		}
	}()

	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	var entries []string
	for _, f := range files {
		fullPath := filepath.Join(sourceName, f.Name())
		entries = append(entries, fullPath)
	}

	return entries, nil
}

func loadNewsFromFile(filePath string) ([]entity.News, error) {
	var news []entity.News
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to read file %v", err)
		return nil, err
	}
	err = json.Unmarshal(jsonData, &news)
	if err != nil {
		log.Printf("Failed to unmarshal JSON data from file %s: %v", filePath, err)
		return nil, err
	}
	return news, nil
}
