package agregator

import (
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/filter"
	"NewsAggregator/cli/internal/parser"
	"fmt"
	"strings"
)

type NewsAggregator struct {
	Sources []string
	Filters []filter.NewsFilter
}

var resources []entity.Resource

func (aggregator NewsAggregator) New() []entity.News {
	var result []entity.News
	initializeResource()
	for _, sourceName := range aggregator.Sources {
		sourceName = strings.TrimSpace(sourceName)
		for _, source := range resources {
			if strings.EqualFold(string(source.Name), sourceName) {
				news, err := parser.New(source.PathToFile).Parse()
				if err != nil {
					fmt.Println(err)
				} else {
					result = append(result, news...)
				}
			}
		}
	}
	for _, newsFilter := range aggregator.Filters {
		result = newsFilter.Filter(result)
	}
	return result
}
func initializeResource() {
	resources = []entity.Resource{
		{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml"},
		{Name: "NBC", PathToFile: "resources/nbc-news.json"},
		{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml"},
		{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"},
		{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html"},
	}
}

//func (agr NewsAggregator) ProcessKeywords(keywords []string) {
//	agr.Filters := append(agr.Filters, &filter.KeywordFilter{Keywords: keywords})
//}
