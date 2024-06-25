package internal

import (
	"news-aggregator/internal/entity"
	"news-aggregator/internal/filters"
	"news-aggregator/internal/sort"
	t "news-aggregator/internal/template"
	"os"
	"strings"
)

type aggregator struct {
	News        map[string][]entity.News
	Sources     string
	NewsFilters []filters.NewsFilter
	SortOptions sort.Options
}

func NewAggregator(news map[string][]entity.News, sources string, newsFilters []filters.NewsFilter, sortParams sort.Options) Aggregate {
	return &aggregator{
		News:        news,
		Sources:     sources,
		NewsFilters: newsFilters,
		SortOptions: sortParams,
	}
}

type Aggregate interface {
	Aggregate() []entity.News
	Print(news []entity.News, keywords string) error
}

// Aggregate aggregates news from the specified Sources and applies NewsFilters.
func (a *aggregator) Aggregate() []entity.News {
	sources := strings.Split(a.Sources, ",")
	news := a.collectNews(sources)
	news = a.applyFilters(news)
	return a.SortOptions.Sort(news)
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
func (a *aggregator) collectNews(sources []string) []entity.News {
	var news []entity.News
	for _, sourceName := range sources {
		sourceName = strings.TrimSpace(sourceName)
		newsFromSource := a.getNewsForSource(sourceName)
		news = append(news, newsFromSource...)
	}
	return news
}

// getNewsForSource fetches news for a single source by comparing it with the list of resources.
func (a *aggregator) getNewsForSource(sourceName string) []entity.News {
	for source, news := range a.News {
		if strings.Contains(source, sourceName) {
			return news
		}
	}
	return nil
}

// applyFilters applies the configured NewsFilters to the aggregated news.
func (a *aggregator) applyFilters(news []entity.News) []entity.News {
	for _, current := range a.NewsFilters {
		news = current.Filter(news)
	}
	return news
}
