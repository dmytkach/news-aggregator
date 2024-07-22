package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

type Fetch struct {
	SourceManager managers.SourceManager
	NewsManager   managers.NewsManager
	FeedManager   managers.FeedManager
}

// UpdateNews from all registered sources and updates the local storage.
func (f Fetch) UpdateNews() error {
	sources, err := f.SourceManager.GetSources()
	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return err
	}
	for _, s := range sources {
		for _, link := range s.PathsToFile {
			err := f.fetchNewsFromSource(entity.Source{
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

// fetchNewsFromSource and updates local storage if the news is not already present.
func (f Fetch) fetchNewsFromSource(resource entity.Source) error {
	for _, path := range resource.PathsToFile {
		news, err := f.FeedManager.FetchFeed(string(path))
		if err != nil {
			log.Printf("Failed to fetch news from %s: %v", path, err)
			continue
		}
		allNews, err := f.NewsManager.GetNewsFromFolder(string(resource.Name))
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
		if len(newsWithoutRepeat) > 0 {
			err = f.NewsManager.AddNews(newsWithoutRepeat, string(news.Name))
			if err != nil {
				log.Printf("Failed to add news for %s: %v", resource.Name, err)
				return err
			}
		}
	}
	return nil
}

func articleExists(existingLinks []entity.Link, newArticle entity.News) bool {
	for _, link := range existingLinks {
		if newArticle.Link == link {
			return true
		}
	}
	return false
}
