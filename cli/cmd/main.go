package main

import (
	"NewsAggregator/cli/internal"
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/filters"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	help := flag.Bool("help", false, "Show all available arguments and their descriptions.")
	sources := flag.String("sources", "", "Select the desired news sources to get the news from. Usage: --sources=bbc,usatoday")
	keywords := flag.String("keywords", "", "Specify the keywords to filter the news by. Usage: --keywords=Ukraine,China")
	dateStart := flag.String("date-start", "", "Specify the start date to filter the news by. Usage: --date-start=2024-05-18")
	dateEnd := flag.String("date-end", "", "Specify the end date to filter the news by. Usage: --date-end=2024-05-19")
	sortOrder := flag.String("sort-order", "DESC", "Specify the sort order for the news items (ASC or DESC). Usage: --sort-order=ASC")
	sortBy := flag.String("sort-by", "date", "Specify the sort criteria for the news items (date or source). Usage: --sort-by=source")

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

	// Sort news based on provided criteria
	internal.SortNews(res, *sortBy, *sortOrder)

	printNews(res, *sortBy, *keywords, sourceList)
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

func printNews(news []entity.News, sortBy, keywords string, sourceList []string) {
	var tmplFile = "cli/internal/entity/news.tmpl"
	funcMap := template.FuncMap{
		"highlight": func(text, keywords string) string {
			for _, keyword := range strings.Split(keywords, ",") {
				text = strings.ReplaceAll(text, keyword, "~~"+keyword+"~~")
			}
			return text
		},
	}
	tmpl, err := template.New("news").Funcs(funcMap).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}

	data := struct {
		Filters  string
		Count    int
		News     []entity.News
		Keywords string
		SortBy   string
	}{
		Filters:  fmt.Sprintf("sources=%s, keywords=%s, sort-by=%s", strings.Join(sourceList, ","), keywords, sortBy),
		Count:    len(news),
		News:     news,
		Keywords: fmt.Sprintf("%s", keywords),
		SortBy:   sortBy,
	}

	err = tmpl.ExecuteTemplate(os.Stdout, "news", data)
	if err != nil {
		panic(err)
	}
}
