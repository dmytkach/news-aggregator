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

	initializeResource()

	sourceList := strings.Split(*sources, ",")
	if len(sourceList) == 0 || (len(sourceList) == 1 && sourceList[0] == "") {
		fmt.Println("Please provide at least one source using the --sources flag.")
		return
	}

	var result []entity.News
	for _, sourceName := range sourceList {
		sourceName = strings.TrimSpace(sourceName)
		for _, source := range resources {
			if strings.EqualFold(string(source.Name), sourceName) {
				news, err := parser.GetParser(source.SourceType).Parse(string(source.PathToFile))
				if err != nil {
					fmt.Println(err)
					continue
				}
				result = append(result, news...)
			}
		}
	}
	newsFilter := filter.NewsFilter{News: &result}
	if *keywords != "" {
		keywordList := strings.Split(*keywords, ",")
		result = newsFilter.FilterByKeywords(keywordList)
	}

	if *dateStart != "" && *dateEnd != "" {
		startDate, err := time.Parse("2006-01-02", *dateStart)
		if err != nil {
			fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
			return
		}
		endDate, err := time.Parse("2006-01-02", *dateEnd)
		if err != nil {
			fmt.Println("Invalid end date format. Please use YYYY-MM-DD.")
			return
		}
		result = newsFilter.FilterByDate(startDate, endDate)
	}

	for _, newsItem := range result {
		fmt.Println(newsItem.ToString())
		fmt.Println("-----------------------")
	}
}

func initializeResource() {
	resources = []entity.Resource{
		{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "NBC", PathToFile: "resources/nbc-news.json", SourceType: "JSON"},
		{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml", SourceType: "RSS"},
		{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html", SourceType: "Html"},
	}
}
