package managers

import (
	"io"
	"log"
	"net/http"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/parser"
	"os"
	"strings"
)

const fileName = "file"

// FeedManager for fetching news feeds.
type FeedManager interface {
	FetchFeed(path string) (entity.Feed, error)
}

// UrlFeed implements the FeedManager for fetching feeds from URLs.
type UrlFeed struct {
}

// FetchFeed downloads and parses the news feed from the given URL.
func (f UrlFeed) FetchFeed(path string) (entity.Feed, error) {
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
	contentType := resp.Header.Get("Content-Type")
	ext := getContentExt(contentType)
	tempFileName := fileName + ext
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		log.Printf("Failed to create temporary file: %v", err)
		return entity.Feed{}, err
	}
	defer func() {
		if err := os.Remove(tempFileName); err != nil {
			log.Printf("Error removing temporary file %s: %v", tempFileName, err)
		}
	}()
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

// getContentExt returns the corresponding file extension based on the content type
func getContentExt(contentType string) string {
	if strings.Contains(contentType, "application/json") {
		return ".json"
	} else if strings.Contains(contentType, "application/rss+xml") ||
		strings.Contains(contentType, "text/xml") {
		return ".xml"
	} else if strings.Contains(contentType, "text/html") {
		return ".html"
	}
	return ""
}
