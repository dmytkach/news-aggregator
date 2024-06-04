package main

import (
	"NewsAggregator/cli/internal"
	"NewsAggregator/cli/internal/filters"
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
		println("Please provide at least one source using the --sources flag.")
		return
	}
	var resFilter []internal.NewsFilter
	resFilter = append(resFilter, processKeywords(*keywords))

	if *dateStart != "" {
		start, err := processDateStart(*dateStart)
		if err != nil || start == nil {
			return
		}
		resFilter = append(resFilter, start)
	}

	if *dateEnd != "" {
		end, err := processDateEnd(*dateEnd)
		if err != nil {
			return
		}
		resFilter = append(resFilter, end)
	}
	res := internal.Aggregate(sourceList, resFilter)
	for _, news := range res {
		fmt.Println(news.ToString())
	}
}

func processDateEnd(dateEnd string) (internal.NewsFilter, error) {
	endDate, err := time.Parse("2006-01-02", dateEnd)
	if err != nil {
		fmt.Println("Invalid end date format. Please use YYYY-MM-DD.")
		return nil, err
	}
	return &filters.DateEnd{EndDate: endDate}, err
}

func processDateStart(dateStart string) (internal.NewsFilter, error) {
	startDate, err := time.Parse("2006-01-02", dateStart)
	if err != nil {
		fmt.Println("Invalid start date format. Please use YYYY-MM-DD.")
		return nil, err
	}
	if startDate.After(time.Now()) {
		println("News for this period is not available.")
		return nil, nil
	}
	return &filters.DateStart{StartDate: startDate}, err
}

func processKeywords(keywords string) internal.NewsFilter {
	keywordList := strings.Split(keywords, ",")
	return &filters.Keyword{Keywords: keywordList}
}
