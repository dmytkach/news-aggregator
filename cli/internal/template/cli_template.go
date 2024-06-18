package template

import (
	"fmt"
	"github.com/wk8/go-ordered-map"
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/sort"
	"strings"
	"text/template"
)

// TemplateData represents the data structure for the news template.
type TemplateData struct {
	Header  Header
	News    []entity.News
	Grouped []*groupedNews
}
type Header struct {
	Sources     string
	Filters     string
	SortOptions sort.Options
}

// groupedNews represents a group of news items.
type groupedNews struct {
	Source   string
	NewsList []entity.News
}

const pathToTemplate = "cli/internal/template/news.tmpl"

func (t TemplateData) Create(keywords string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"highlight": func(text string) string {
			if len(keywords) == 0 {
				return text
			}
			for _, keyword := range strings.Split(strings.ToLower(keywords), ",") {
				text = strings.ReplaceAll(text, keyword, "~~"+keyword+"~~")
			}
			return text
		},
		"toString": func(v interface{}) string {
			return fmt.Sprintf("%v", v)
		},
	}

	tmpl, err := template.New("news").Funcs(funcMap).ParseFiles(pathToTemplate)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return tmpl, err
}

// Prepare the template data for rendering.
func (t TemplateData) Prepare() TemplateData {
	groupedMap := group(t.News)
	var groupedList []*groupedNews
	for el := groupedMap.Oldest(); el != nil; el = el.Next() {
		source := el.Key.(string)
		newsList := el.Value.([]entity.News)
		groupedList = append(groupedList, &groupedNews{Source: source, NewsList: newsList})
	}
	t.Grouped = groupedList
	return t
}

// group the news items by their source.
func group(news []entity.News) *orderedmap.OrderedMap {
	grouped := orderedmap.New()
	for _, item := range news {
		if value, ok := grouped.Get(item.Source); ok {
			grouped.Set(item.Source, append(value.([]entity.News), item))
		} else {
			grouped.Set(item.Source, []entity.News{item})
		}
	}
	return grouped
}
