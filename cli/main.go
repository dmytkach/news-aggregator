package main

import (
	"flag"
	"news-aggregator/internal"
	"news-aggregator/internal/filters"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
)

// main is the entry point of the news-aggregator CLI application.
func main() {
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Usage: --sources=bbc,usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by. Usage: --keywords=Ukraine,China")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by. Usage: --date-start=2024-05-18")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by. Usage: --date-end=2024-05-19")
	sortOrder := flag.String("sort-order", "ASC", "Specify the sort order for the news items (ASC or DESC). The default is ASC. Usage: --sort-order=ASC")
	sortBy := flag.String("sort-by", "source", "Specify the sort criteria for the news items (date or source). The default is source. Usage: --sort-by=source")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	sortOptions := sort.Options{
		Criterion: *sortBy,
		Order:     *sortOrder,
	}
	resources := initializeDefaultResource()
	config := validator.Config{
		Sources:          *sources,
		AvailableSources: resources,
		DateStart:        *dateStart,
		DateEnd:          *dateEnd,
		SortOptions:      sortOptions,
	}

	v := validator.NewValidator(config)
	if !v.Validate() {
		return
	}
	a := internal.NewAggregator(
		resources,
		*sources,
		filters.InitializeFilters(keywords, dateStart, dateEnd),
		sortOptions)
	news, err := a.Aggregate()
	if err != nil {
		print(err)
		return
	}
	err = a.Print(news, *keywords)
	if err != nil {
		print(err)
		return
	}
}
func initializeDefaultResource() map[string]string {
	return map[string]string{
		"bbc":        "resources/bbc-world-category-19-05-24.xml",
		"nbc":        "resources/nbc-news.json",
		"abc":        "resources/abcnews-international-category-19-05-24.xml",
		"washington": "resources/washingtontimes-world-category-19-05-24.xml",
		"usatoday":   "resources/usatoday-world-news.html",
	}
}
