package internal

import (
	"NewsAggregator/cli/internal/entity"
	"strings"
)

var resources = initializeResource()

// Aggregate aggregates news from the specified sources and applies filters.
func Aggregate(sources []string, filters []NewsFilter) []entity.News {
	sourceNews := collectNews(sources)
	return applyFilters(sourceNews, filters)
}

// collectNews collects news from all specified resources.
func collectNews(sources []string) []entity.News {
	var result []entity.News
	for _, sourceName := range sources {
		sourceName = strings.TrimSpace(sourceName)
		newsFromSource := getForSource(sourceName)
		result = append(result, newsFromSource...)
	}
	return result
}

// getForSource fetches news for a single source by comparing it with the list of resources.
func getForSource(sourceName string) []entity.News {
	var result []entity.News
	for _, resource := range resources {
		if strings.EqualFold(string(resource.Name), sourceName) {
			news := getResourceNews(resource)
			result = append(result, news...)
		}
	}
	if len(result) == 0 {
		print("Error news source: " + sourceName + ". Available news sources: ")
		for _, resource := range resources {
			print("  " + resource.Name)
		}
		print("\n")
	}
	return result
}

// getResourceNews parses news from a single resource.
func getResourceNews(resource entity.Resource) []entity.News {
	news, err := New(resource.PathToFile).Parse()
	if err != nil {
		print("Error fetching news from source: " + string(resource.Name))
	}
	for i := range news {
		news[i].Source = string(resource.Name)
	}
	return news
}

// applyFilters applies the configured filters to the aggregated news.
func applyFilters(news []entity.News, filters []NewsFilter) []entity.News {
	for _, current := range filters {
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
