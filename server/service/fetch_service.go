package service

import (
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
)

type FetchService struct {
	SourceRepo managers.SourceManager
	NewsRepo   managers.NewsManager
}

// UpdateNews from all registered sources and updates the local storage.
func (fetcher FetchService) UpdateNews() error {
	sources, err := fetcher.SourceRepo.GetSources()
	if err != nil {
		log.Printf("Error fetching sources: %v", err)
		return err
	}
	for _, s := range sources {
		for _, link := range s.PathsToFile {
			err := fetcher.fetchNewsFromSource(entity.Source{
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
func (fetcher FetchService) fetchNewsFromSource(resource entity.Source) error {
	for _, path := range resource.PathsToFile {
		news, err := managers.FetchFeed(path)
		if err != nil {
			log.Printf("Failed to fetch news from %s: %v", path, err)
			continue
		}
		allNews, err := fetcher.NewsRepo.GetNewsFromFolder(string(resource.Name))
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
		err = fetcher.NewsRepo.AddNews(newsWithoutRepeat, string(news.Name))
		if err != nil {
			log.Printf("Failed to add news for %s: %v", resource.Name, err)
			return err
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
