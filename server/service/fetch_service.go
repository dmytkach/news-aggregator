package service

import (
	"io"
	"log"
	"net/http"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
	"os"
)

const tempFileName = "tempfile.xml"

// FetchNews from all registered sources and updates the local storage.
func FetchNews() error {
	sources, err := managers.GetSources()
	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return err
	}
	for _, s := range sources {
		links, err := managers.GetSource(string(s.Name))
		if err != nil {
			log.Printf("Error fetching feeds for source %s: %v", s, err)
			return err
		}
		for _, link := range links.PathsToFile {
			err := fetchNewsFromResource(entity.Source{
				Name:        s.Name,
				PathsToFile: []entity.PathToFile{link},
			})
			if err != nil {
				log.Printf("Error fetching news from resource %s: %v", s, err)
				return err
			}
		}
	}
	return nil
}

// fetchNewsFromResource and updates local storage if the news is not already present.
func fetchNewsFromResource(resource entity.Source) error {
	for _, path := range resource.PathsToFile {
		news, err := fetchFeed(path)
		if err != nil {
			log.Printf("Failed to fetch news from %s: %v", path, err)
			continue
		}
		allNews, err := managers.GetNewsFromFolder(string(resource.Name))
		if err != nil {
			log.Printf("Failed to get existing news for %s: %v", resource.Name, err)
			return err
		}
		allNewsLink := make([]entity.Link, 0)
		for _, n := range allNews {
			allNewsLink = append(allNewsLink, n.Link)
		}

		newsWithoutRepeat := make([]entity.News, 0)
		for _, loadedNews := range news.News {
			if !articleExists(allNewsLink, loadedNews) {
				newsWithoutRepeat = append(newsWithoutRepeat, loadedNews)
			}
		}
		err = managers.AddNews(newsWithoutRepeat, string(news.Name))
		if err != nil {
			log.Printf("Failed to add news for %s: %v", resource.Name, err)
			return err
		}
	}
	return nil
}

// fetchFeed downloads and parses the news feed from the given URL.
func fetchFeed(url entity.PathToFile) (entity.Feed, error) {
	resp, err := http.Get(string(url))
	if err != nil {
		log.Println("Failed to download feed", http.StatusInternalServerError)
		return entity.Feed{}, err
	}
	defer resp.Body.Close()
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

	feed, err := managers.GetNewsFromFile(tempFileName)
	if err != nil {
		log.Printf("Failed to parse feed from file: %v", err)
		return entity.Feed{}, err
	}
	return feed, nil
}

func articleExists(existingLinks []entity.Link, newArticle entity.News) bool {
	for _, link := range existingLinks {
		if newArticle.Link == link {
			return true
		}
	}
	return false
}
