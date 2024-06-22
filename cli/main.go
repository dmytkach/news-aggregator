package main

import (
	"flag"
	"news-aggregator/internal"
	"news-aggregator/internal/sort"
	"news-aggregator/internal/validator"
	"news-aggregator/storage"
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
	s := storage.NewMemoryStorage()
	initializeDefaultResource(s)
	v := validator.Validator{
		Sources:          *sources,
		AvailableSources: s.AvailableSources(),
		DateStart:        *dateStart,
		DateEnd:          *dateEnd,
	}
	if !v.Validate() {
		return
	}
	a := internal.NewAggregator(
		s.GetAll(),
		*sources,
		internal.InitializeFilters(keywords, dateStart, dateEnd),
		sort.Options{
			Criterion: *sortBy,
			Order:     *sortOrder,
		})
	news := a.Aggregate()
	err := a.Print(news, *keywords)
	if err != nil {
		print(err)
	}
}
func initializeDefaultResource(s storage.Storage) {
	s.Set("bbc", "resources/bbc-world-category-19-05-24.xml")
	s.Set("nbc", "resources/nbc-news.json")
	s.Set("abc", "resources/abcnews-international-category-19-05-24.xml")
	s.Set("washington", "resources/washingtontimes-world-category-19-05-24.xml")
	s.Set("usatoday", "resources/usatoday-world-news.html")

}
