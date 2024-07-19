package managers

import (
	"io"
	"log"
	"net/http"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/parser"
	"os"
)

const tempFileName = "tempfile.xml"

type FeedManager interface {
	Fetch(path string) (entity.Feed, error)
}
type UrlFeed struct {
}

// Fetch downloads and parses the news feed from the given URL.
func (f UrlFeed) Fetch(path string) (entity.Feed, error) {
	resp, err := http.Get(path)
	if err != nil {
		log.Println("Failed to download feed", http.StatusInternalServerError)
		return entity.Feed{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
		return entity.Feed{}, err
	}
	defer os.Remove(tempFileName)

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		log.Printf("Failed to write response to file: %v", err)
		return entity.Feed{}, err
	}
	err = tempFile.Close()
	if err != nil {
		log.Printf("Failed to close temporary file: %v", err)
		return entity.Feed{}, err
	}

	feed, err := getFeedFromFile(tempFileName)
	if err != nil {
		log.Printf("Failed to parse feed from file: %v", err)
		return entity.Feed{}, err
	}
	return feed, nil
}

// getFeedFromFile using parsers.
func getFeedFromFile(filePath string) (entity.Feed, error) {
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
