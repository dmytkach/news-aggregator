package main

import (
	"NewsAggregator/cli/internal/aggregator"
	"NewsAggregator/cli/internal/filter"
	"flag"
	"fmt"
	"strings"
	"time"
)

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

	sourceList := strings.Split(*sources, ",")
	if len(sourceList) == 0 {
		fmt.Println("Please provide at least one source using the --sources flag.")
		return
	}
	var filters []aggregator.NewsFilter
	filters = append(filters, processKeywords(*keywords))

	if *dateStart != "" {
		start, err := processDateStart(*dateStart)
		if err != nil {
			return
		}
		filters = append(filters, start)
	}

	if *dateEnd != "" {
		end, err := processDateEnd(*dateEnd)
		if err != nil {
			return
		}
		filters = append(filters, end)
	}
	res := aggregator.NewsAggregator{Sources: sourceList, Filters: filters}.New()
	for _, news := range res {
		fmt.Println(news.ToString())
	}
}

func processDateEnd(dateEnd string) (aggregator.NewsFilter, error) {
	endDate, err := time.Parse("2006-01-02", dateEnd)
	if err != nil {
		fmt.Println("Invalid end date format. Please use YYYY-MM-DD.")
		return nil, err
	}
	return &filter.DateEndFilter{EndDate: endDate}, err
}

func processDateStart(dateStart string) (aggregator.NewsFilter, error) {
	startDate, err := time.Parse("2006-01-02", dateStart)
	if err != nil {
		fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
		return nil, err
	}
	return &filter.DateStartFilter{StartDate: startDate}, err
}

func processKeywords(keywords string) aggregator.NewsFilter {
	keywordList := strings.Split(keywords, ",")
	return &filter.KeywordFilter{Keywords: keywordList}
}
