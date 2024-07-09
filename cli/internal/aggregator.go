package internal

import (
	"news-aggregator/internal/entity"
	"news-aggregator/internal/sort"
	t "news-aggregator/internal/template"
	"os"
	"strings"
)

// Parser provides an API for a news parser capable of processing a specific file type.
type Parser interface {
	CanParseFileType(ext string) bool
	Parse() ([]entity.News, error)
}

// aggregator aggregates news data from various sources.
type aggregator struct {
	Resources   map[string]string
	Sources     string
	NewsFilters []NewsFilter
	SortOptions sort.Options
}

// NewAggregator creates a new instance of an aggregator with the given resources, sources,
// news filters, and sorting options.
func NewAggregator(resources map[string]string, sources string, newsFilters []NewsFilter, sortParams sort.Options) Aggregate {
	return &aggregator{
		Resources:   resources,
		Sources:     sources,
		NewsFilters: newsFilters,
		SortOptions: sortParams,
	}
}

// Aggregate provides API for aggregating news data from various sources
// and for printing the aggregated news using predefined templates.
type Aggregate interface {
	Aggregate() ([]entity.News, error)
	Print(news []entity.News, keywords string) error
}

// Aggregate aggregates news from the specified Sources and applies NewsFilters.
func (a *aggregator) Aggregate() ([]entity.News, error) {
	sources := strings.Split(a.Sources, ",")
	news, err := a.collectNews(sources)
	if err != nil {
		return nil, err
	}
	news = a.applyFilters(news)
	return a.SortOptions.Sort(news), nil
}

// Print news according to the created template
func (a *aggregator) Print(news []entity.News, keywords string) error {
	template := t.Data{
		News: news,
		Header: t.Header{
			Sources:     a.Sources,
			SortOptions: a.SortOptions,
		},
	}
	if len(a.NewsFilters) != 0 {
		var filtersInfo string
		for i := range a.NewsFilters {
			filtersInfo += a.NewsFilters[i].String() + " "
		}
		filtersInfo = " filters:" + filtersInfo
		template.Header.Filters = filtersInfo
	}
	tmpl, err := template.Create(keywords)
	if err != nil {
		return err
	}
	data := template.Prepare()
	err = tmpl.ExecuteTemplate(os.Stdout, "news", data)
	if err != nil {
		return err
	}
	return nil
}

// collectNews collects news from all specified resources.
func (a *aggregator) collectNews(sources []string) ([]entity.News, error) {
	var news []entity.News
	for _, sourceName := range sources {
		sourceName = strings.TrimSpace(sourceName)
		newsFromSource, err := a.getNewsForSource(sourceName)
		if err != nil {
			return nil, err
		}
		news = append(news, newsFromSource...)
	}
	return news, nil
}

// getNewsForSource fetches news for a single source by comparing it with the list of resources.
func (a *aggregator) getNewsForSource(sourceName string) ([]entity.News, error) {
	value, _ := a.Resources[strings.ToLower(sourceName)]
	resourceNews, err := a.getResourceNews(entity.PathToFile(value))
	if err != nil {
		println("File path specified incorrectly for source " + sourceName)
		return nil, err
	}
	for i := range resourceNews {
		resourceNews[i].Source = sourceName
	}
	return resourceNews, nil
}

// getResourceNews parses news from a single resource.
func (a *aggregator) getResourceNews(path entity.PathToFile) ([]entity.News, error) {
	p, err := GetFileParser(path)
	if err != nil {
		return nil, err
	}
	news, err := p.Parse()
	if err != nil {
		return nil, err
	}
	return news, nil
}

// applyFilters applies the configured NewsFilters to the aggregated news.
func (a *aggregator) applyFilters(news []entity.News) []entity.News {
	for _, current := range a.NewsFilters {
		news = current.Filter(news)
	}
	return news
}
