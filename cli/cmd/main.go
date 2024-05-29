package main

import (
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/filter"
	"NewsAggregator/cli/internal/parser"
	"flag"
	"fmt"
	"strings"
	"time"
)

var resources []entity.Resource

func main() {
	resources = append(resources, entity.Resource{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml"})
	resources = append(resources, entity.Resource{Name: "NBC", PathToFile: "resources/nbc-news.json"})
	resources = append(resources, entity.Resource{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml"})
	resources = append(resources, entity.Resource{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"})
	resources = append(resources, entity.Resource{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html"})

	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Usage: --sources=bbc,usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by. Usage: --keywords=Ukraine,China")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by. Usage: --date-start=2024-05-18")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by. Usage: --date-end=2024-05-19")

	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	sourceList := strings.Split(*sources, ",")
	if len(sourceList) == 0 {
		fmt.Println("Please provide at least one source using the --sources flag.")
		return
	}

	var result []entity.News
	for _, sourceName := range sourceList {
		sourceName = strings.TrimSpace(sourceName)
		for _, source := range resources {
			if strings.EqualFold(string(source.Name), sourceName) {
				news, err := parser.GetParser(source.PathToFile).Parse()
				if err != nil {
					fmt.Println(err)
				} else {
					result = append(result, news...)
				}
			}
		}
	}

	var filters []filter.NewsFilter
	if *keywords != "" {
		keywordList := strings.Split(*keywords, ",")
		filters = append(filters, &filter.KeywordFilter{Keywords: keywordList})
	}

	if *dateStart != "" {
		startDate, err := time.Parse("2006-01-02", *dateStart)
		if err != nil {
			fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
			return
		}
		filters = append(filters, &filter.DateStartFilter{StartDate: startDate})
	}

	if *dateEnd != "" {
		endDate, err := time.Parse("2006-01-02", *dateEnd)
		if err != nil {
			fmt.Println("Invalid end date format. Please use YYYY-MM-DD.")
			return
		}
		filters = append(filters, &filter.DateEndFilter{EndDate: endDate})
	}

	for _, newsFilter := range filters {
		result = newsFilter.Filter(result)
	}

	for _, newsItem := range result {
		fmt.Println(newsItem.ToString())
	}
}

//
//func initializeResource() {
//	resources = []entity.Resource{
//		{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml"},
//		{Name: "NBC", PathToFile: "resources/nbc-news.json"},
//		{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml"},
//		{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"},
//		{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html"},
//	}
//}
