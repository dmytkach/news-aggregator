package main

import (
	"flag"
	"news-aggregator/internal"
	"news-aggregator/internal/filters"
	"news-aggregator/internal/validator"
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
	sortOrder := flag.String("sort-order", "ASC", "Specify the sort order for the news items (ASC or DESC). The default is ASC. Usage: --sort-order=ASC")
	sortBy := flag.String("sort-by", "source", "Specify the sort criteria for the news items (date or source). The default is source. Usage: --sort-by=source")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}
	resources := initializeDefaultResource()
	v := validator.Validator{
		Sources:          *sources,
		AvailableSources: resources,
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
		startDate, _ := time.Parse(validator.DateFormat, *dateStart)
		newsFilters = append(newsFilters, &filters.DateStart{StartDate: startDate})
	}
	if len(*dateEnd) > 0 {
		endDate, _ := time.Parse(validator.DateFormat, *dateEnd)
		newsFilters = append(newsFilters, &filters.DateEnd{EndDate: endDate})
	}
	a := internal.NewAggregator(resources, sourceList, newsFilters)
	news := a.Aggregate()
	sort := internal.NewSort(*sortBy, *sortOrder)
	news = sort.Apply(news)
	internal.Template{News: news, Criterion: *sortBy}.Apply(newsFilters, *sortOrder, *keywords)
}

func initializeDefaultResource() map[string]string {
	return map[string]string{
		"bbc":        "resources/bbc-world-category-19-05-24.xml",
		"nbc":        "resources/nbc-news.json",
		"abc":        "resources/abcnews-international-category-19-05-24.xml",
		"washington": "resources/washingtontimes-world-category-19-05-24.xml",
		"usatoday":   "resources/usatoday-world-news.html",
	}
	//return []entity.Resource{
	//	{Name: "BBC", PathToFile: "resources/bbc-world-category-19-05-24.xml"},
	//	{Name: "NBC", PathToFile: "resources/nbc-news.json"},
	//	{Name: "ABC", PathToFile: "resources/abcnews-international-category-19-05-24.xml"},
	//	{Name: "Washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"},
	//	{Name: "USAToday", PathToFile: "resources/usatoday-world-news.html"},
	//}
}
