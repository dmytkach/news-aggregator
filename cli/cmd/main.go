package main

import (
	"flag"
	"news-aggregator/cli/internal"
	"news-aggregator/cli/internal/entity"
	"news-aggregator/cli/internal/filters"
	"news-aggregator/cli/internal/validator"
	"strings"
	"time"
)

// main is the entry point of the news-aggregator CLI application.
func main() {
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Usage: --sources=bbc,usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by. Usage: --keywords=Ukraine,China")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by. Usage: --date-start=2024-05-18")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by. Usage: --date-end=2024-05-19")
	sortOrder := flag.String("sort-order", "ASC", "Specify the sort order for the news items (ASC or DESC). Usage: --sort-order=ASC")
	sortBy := flag.String("sort-by", "source", "Specify the sort criteria for the news items (date or source). Usage: --sort-by=source")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	resources := initializeResource()
	v := validator.Validator{
		Sources:          *sources,
		AvailableSources: entity.GetResourceNames(resources),
		DateStart:        *dateStart,
		DateEnd:          *dateEnd,
	}
	if !v.Validate() {
		return
	}
	sourceList := strings.Split(v.Sources, ",")
	var newsFilters []internal.NewsFilter
	if len(*keywords) > 0 {
		keywordList := strings.Split(*keywords, ",")
		newsFilters = append(newsFilters, &filters.Keyword{Keywords: keywordList})
	}
	if len(*dateStart) > 0 {
		startDate, _ := time.Parse("2006-01-02", *dateStart)
		newsFilters = append(newsFilters, &filters.DateStart{StartDate: startDate})
	}
	if len(*dateEnd) > 0 {
		endDate, _ := time.Parse("2006-01-02", *dateEnd)
		newsFilters = append(newsFilters, &filters.DateEnd{EndDate: endDate})
	}
	aggregator := internal.Aggregator{Resources: resources, Sources: sourceList, NewsFilters: newsFilters}
	news := aggregator.Aggregate()
	news = internal.Sort(news, *sortBy, *sortOrder)
	Template{News: news, Criterion: *sortBy}.apply(newsFilters, *sortOrder, *keywords)
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
