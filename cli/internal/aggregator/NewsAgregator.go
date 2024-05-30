package aggregator

import (
	"NewsAggregator/cli/internal/entity"
	"fmt"
	"strings"
)

// NewsAggregator is responsible for aggregating and filtering news articles
// from various sources specified in the Sources field and applying filters
// specified in the Filters field.
type NewsAggregator struct {
	Sources []string
	Filters []NewsFilter
}

// New aggregated news from these sources Using these filters.
// It initializes the list of resources, analyzes news from each source and applies
// filters to return a filtered set of news articles.
func (aggregator NewsAggregator) New() []entity.News {
	var result []entity.News
	resources := initializeResource()
	for _, sourceName := range aggregator.Sources {
		sourceName = strings.TrimSpace(sourceName)
		for _, source := range resources {
			if strings.EqualFold(string(source.Name), sourceName) {
				news, err := New(source.PathToFile).Parse()
				if err != nil {
					fmt.Println(err)
				} else {
					result = append(result, news...)
				}
			}
		}
	}
	for _, news := range aggregator.Filters {
		result = news.Filter(result)
	}
	return result
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
