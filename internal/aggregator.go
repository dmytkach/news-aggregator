package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/filters"
	"news-aggregator/internal/sort"
	t "news-aggregator/internal/template"
	"os"
	"strings"
)

type aggregator struct {
	Resources   map[string][]string
	Sources     string
	NewsFilters []filters.NewsFilter
	SortOptions sort.Options
}

func NewAggregator(news map[string][]string, sources string, newsFilters []filters.NewsFilter, sortParams sort.Options) Aggregate {
	return &aggregator{
		Resources:   news,
		Sources:     sources,
		NewsFilters: newsFilters,
		SortOptions: sortParams,
	}
}

type Aggregate interface {
	Aggregate() ([]entity.News, error)
	Print(news []entity.News, keywords string) error
}

// Aggregate news from the specified Sources and applies NewsFilters.
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
		log.Printf("Error creating template: %v", err)
		return err
	}
	data := template.Prepare()
	err = tmpl.ExecuteTemplate(os.Stdout, "news", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		return err
	}
	return nil
}

// collectNews from all specified resources.
func (a *aggregator) collectNews(sources []string) ([]entity.News, error) {
	var news []entity.News
	for _, sourceName := range sources {
		sourceName = strings.ToLower(strings.TrimSpace(sourceName))
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
	var news []entity.News
	for source, path := range a.Resources {
		if strings.Contains(source, sourceName) {
			for _, b := range path {
				newsFromResource, err := a.getResourceNews(entity.PathToFile(b))
				if err != nil {
					return nil, err
				}
				news = append(news, newsFromResource...)
			}
		}
	}
	return news, nil
}

// getResourceNews from a single resource.
func (a *aggregator) getResourceNews(path entity.PathToFile) ([]entity.News, error) {
	file, err := os.Open(string(path))
	if err != nil {
		log.Printf("Failed to open file %s: %v", path, err)
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil && err == nil {
			err = fmt.Errorf("error closing file: %w", closeErr)
		}
	}(file)
	var articles []entity.News
	if err := json.NewDecoder(file).Decode(&articles); err != nil {
		log.Printf("Error decoding file %s: %v", path, err)
		return nil, err
	}
	return articles, nil
}

// applyFilters applies the configured NewsFilters to the aggregated news.
func (a *aggregator) applyFilters(news []entity.News) []entity.News {
	for _, current := range a.NewsFilters {
		news = current.Filter(news)
	}
	return news
}
