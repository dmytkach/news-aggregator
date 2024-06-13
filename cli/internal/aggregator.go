package internal

import (
	"news-aggregator/cli/internal/entity"
	"strings"
)

type Aggregator struct {
	Resources   []entity.Resource
	Sources     []string
	NewsFilters []NewsFilter
}

// Aggregate aggregates news from the specified Sources and applies NewsFilters.
func (a *Aggregator) Aggregate() []entity.News {
	news := a.collectNews(a.Sources)
	return a.applyFilters(news)
}

// collectNews collects news from all specified resources.
func (a *Aggregator) collectNews(sources []string) []entity.News {
	var news []entity.News
	for _, sourceName := range sources {
		sourceName = strings.TrimSpace(sourceName)
		newsFromSource := a.getNewsForSource(sourceName)
		news = append(news, newsFromSource...)
	}
	return news
}

// getNewsForSource fetches news for a single source by comparing it with the list of resources.
func (a *Aggregator) getNewsForSource(sourceName string) []entity.News {
	var news []entity.News
	for _, resource := range a.Resources {
		if strings.EqualFold(string(resource.Name), sourceName) {
			resourceNews := a.getResourceNews(resource.PathToFile)
			if resourceNews != nil {
				news = append(news, resourceNews...)
			}
		}
	}
	if len(news) == 0 {
		return nil
	}
	for i := range news {
		news[i].Source = sourceName
	}
	return news
}

// getResourceNews parses news from a single resource.
func (a *Aggregator) getResourceNews(path entity.PathToFile) []entity.News {
	p := New(path)
	if p == nil {
		print("Error news source format: ")
		return nil
	}
	news, err := p.Parse()
	if err != nil {
		print("Error parse news from source")
		return nil
	}
	return news
}

// applyFilters applies the configured NewsFilters to the aggregated news.
func (a *Aggregator) applyFilters(news []entity.News) []entity.News {
	for _, current := range a.NewsFilters {
		news = current.Filter(news)
	}
	return news
}
