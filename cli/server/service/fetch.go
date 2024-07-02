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
	resources, err := managers.GetAllSourcesNames()
	if err != nil {
		return err
	}
	for _, name := range resources {
		links, err := managers.GetSourcesFeeds(name)
		if err != nil {
			return err
		}
		for _, link := range links {
			err := fetchNewsFromResource(entity.Resource{
				Name:       entity.ResourceName(name),
				PathToFile: entity.PathToFile(link),
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// fetchNewsFromResource and updates local storage if the news is not already present.
func fetchNewsFromResource(resource entity.Resource) error {
	news, err := fetchNewsFromResponse(resource.PathToFile)
	if err != nil {
		print("Failed to parse feed", http.StatusInternalServerError)
		return err
	}
	allNews, err := managers.GetNewsFromFolder(string(resource.Name))
	if err != nil {
		print("Failed to parse static resources", http.StatusInternalServerError)
		return err
	}
	allNewsLink := make([]entity.Link, 0)
	for _, n := range allNews {
		allNewsLink = append(allNewsLink, n.Link)
	}

	newsWithoutRepeat := make([]entity.News, 0)
	for _, loadedNews := range news {
		if !articleExists(allNewsLink, loadedNews) {
			newsWithoutRepeat = append(newsWithoutRepeat, loadedNews)
		}
	}
	err = managers.AddNews(newsWithoutRepeat, string(resource.Name))
	if err != nil {
		return err
	}
	return nil
}

// fetchNewsFromResponse downloads and parses the news feed from the given URL.
func fetchNewsFromResponse(url entity.PathToFile) ([]entity.News, error) {
	resp, err := http.Get(string(url))
	if err != nil {
		log.Println("Failed to download feed", http.StatusInternalServerError)
		return nil, err
	}
	defer resp.Body.Close()
	tempFile, err := os.Create(tempFileName)
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFileName)

	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return nil, err
	}
	err = tempFile.Close()
	if err != nil {
		return nil, err
	}

	news, err := managers.GetNewsFromFile(tempFileName)
	if err != nil {
		return nil, err
	}
	return news, nil
}
func articleExists(existingLinks []entity.Link, newArticle entity.News) bool {
	for _, link := range existingLinks {
		if newArticle.Link == link {
			return true
		}
	}
	return false
}
