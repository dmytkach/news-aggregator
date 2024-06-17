package internal

import (
	"news-aggregator/internal/entity"
	"os"
	"strings"
)

type aggregator struct {
	Resources   map[string]string
	Sources     string
	NewsFilters []NewsFilter
	SortOptions SortOptions
}

func NewAggregator(resources map[string]string, sources string, newsFilters []NewsFilter, sortParams SortOptions) Aggregate {
	return &aggregator{
		Resources:   resources,
		Sources:     sources,
		NewsFilters: newsFilters,
		SortOptions: sortParams,
	}
}

type Aggregate interface {
	Aggregate() []entity.News
	Print(news []entity.News, keywords string)
}

// Aggregate aggregates news from the specified Sources and applies NewsFilters.
func (a *aggregator) Aggregate() []entity.News {
	sources := strings.Split(a.Sources, ",")
	news := a.collectNews(sources)
	news = a.applyFilters(news)
	return a.SortOptions.Apply(news)
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

// Print the results.
func (a *aggregator) Print(news []entity.News, keywords string) {
	t := Template{
		News: news,
		Info: TemplateInfo{
			Sources:     a.Sources,
			SortOptions: a.SortOptions,
		},
	}
	tmpl := t.CreateTemplate(keywords)
	tmpl1 := t.Prepare(a.NewsFilters)
	err := tmpl.ExecuteTemplate(os.Stdout, "news", tmpl1)
	if err != nil {
		panic(err)
	}
}

// getNewsForSource fetches news for a single source by comparing it with the list of resources.
func (a *aggregator) getNewsForSource(sourceName string) []entity.News {
	var news []entity.News
	value, _ := a.Resources[strings.ToLower(sourceName)]
	resourceNews, err := a.getResourceNews(entity.PathToFile(value))
	if err != nil && resourceNews != nil {
		print(err)
	}
	news = append(news, resourceNews...)
	if len(news) == 0 {
		return nil
	}
	for i := range news {
		news[i].Source = sourceName
	}
	return news
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
