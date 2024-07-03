package main

import (
	"flag"
	"log"
	"news-aggregator/internal"
	"news-aggregator/internal/initializers"
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
	resources, err := initializers.LoadStaticResourcesFromFolder("server-news/")
	if err != nil {
		return
	}
	availableSources := make([]string, 0)
	for sourceName := range resources {
		availableSources = append(availableSources, sourceName)
	}
	v := validator.Validator{
		Sources:          *sources,
		AvailableSources: availableSources,
		DateStart:        *dateStart,
		DateEnd:          *dateEnd,
	}
	if !v.Validate() {
		log.Println("Invalid parameters")
		return
	}
	a := internal.NewAggregator(
		resources,
		*sources,
		initializers.InitializeFilters(keywords, dateStart, dateEnd),
		sort.Options{
			Criterion: *sortBy,
			Order:     *sortOrder,
		})
	news, err := a.Aggregate()
	if err != nil {
		log.Println(err)
		return
	}
	err = a.Print(news, *keywords)
	if err != nil {
		log.Println(err)
		return
	}
}
