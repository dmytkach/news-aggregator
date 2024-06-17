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
	a := internal.NewAggregator(
		resources,
		*sources,
		initializeFilters(keywords, dateStart, dateEnd),
		internal.SortOptions{
			Criterion: *sortBy,
			Order:     *sortOrder,
		})
	news := a.Aggregate()
	a.Print(news, *keywords)
}

// initializeFilters based on provided parameters.
func initializeFilters(keywords, dateStart, dateEnd *string) []internal.NewsFilter {
	var newsFilters []internal.NewsFilter

	if keywordFilter := convertKeywords(keywords); keywordFilter != nil {
		newsFilters = append(newsFilters, keywordFilter)
	}
	if dateStartFilter := convertDateStart(dateStart); dateStartFilter != nil {
		newsFilters = append(newsFilters, dateStartFilter)
	}
	if dateEndFilter := convertDateEnd(dateEnd); dateEndFilter != nil {
		newsFilters = append(newsFilters, dateEndFilter)
	}

	return newsFilters
}

// convertKeywords a comma-separated keyword string into a filter.
func convertKeywords(keywords *string) *filters.Keyword {
	if len(*keywords) > 0 {
		keywordList := strings.Split(*keywords, ",")
		return &filters.Keyword{Keywords: keywordList}
	}
	return nil
}

// convertDateStart string into a DateStart filter.
func convertDateStart(dateStart *string) *filters.DateStart {
	if len(*dateStart) > 0 {
		startDate, _ := time.Parse(validator.DateFormat, *dateStart)
		return &filters.DateStart{StartDate: startDate}
	}
	return nil
}

// convertDateEnd string into a DateEnd filter.
func convertDateEnd(dateEnd *string) *filters.DateEnd {
	if len(*dateEnd) > 0 {
		endDate, _ := time.Parse(validator.DateFormat, *dateEnd)
		return &filters.DateEnd{EndDate: endDate}
	}
	return nil
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
