package internal

import (
	"NewsAggregator/cli/internal/entity"
	"strings"
)

// NewsAggregator is responsible for aggregating and filtering news articles
// from various sources specified in the Sources field and applying filters
// specified in the Filters field.
type NewsAggregator struct {
	Sources []string
	Filters []NewsFilter
}

// New aggregates news from the specified sources and applies filters.
func (aggregator NewsAggregator) New() []entity.News {
	resources := initializeResource()
	sourceNews := aggregator.collectNews(resources)
	return aggregator.applyFilters(sourceNews)
}

// collectNews collects news from all specified resources.
func (aggregator NewsAggregator) collectNews(resources []entity.Resource) []entity.News {
	var result []entity.News
	for _, sourceName := range aggregator.Sources {
		sourceName = strings.TrimSpace(sourceName)
		newsFromSource := aggregator.getForSource(sourceName, resources)
		result = append(result, newsFromSource...)
	}
	return result
}

// getForSource fetches news for a single source by comparing it with the list of resources.
func (aggregator NewsAggregator) getForSource(sourceName string, resources []entity.Resource) []entity.News {
	var result []entity.News
	for _, resource := range resources {
		if strings.EqualFold(string(resource.Name), sourceName) {
			news := aggregator.getResourceNews(resource)
			result = append(result, news...)
		}
	}
	return result
}

// getResourceNews parses news from a single resource.
func (aggregator NewsAggregator) getResourceNews(resource entity.Resource) []entity.News {
	news, err := New(resource.PathToFile).Parse()
	if err != nil {
		print("Error fetching news from source: " + string(resource.Name))
	}
	return news
}

// applyFilters applies the configured filters to the aggregated news.
func (aggregator NewsAggregator) applyFilters(news []entity.News) []entity.News {
	for _, current := range aggregator.Filters {
		news = current.Filter(news)
	}
	return news
}
func initializeResource() []entity.Resource {
	return []entity.Resource{
		{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml"},
		{Name: "NBC", PathToFile: "resources/nbc-news.json"},
		{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml"},
		{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"},
		{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html"},
	}
}
