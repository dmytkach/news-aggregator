package internal

import (
	"fmt"
	"github.com/wk8/go-ordered-map"
	"log"
	"news-aggregator/internal/entity"
	"strings"
	"text/template"
)

// Template represents the data structure for the news template.
type Template struct {
	Filters   string
	News      []entity.News
	Criterion string
	Grouped   []*groupedNews
}

// groupedNews represents a group of news items.
type groupedNews struct {
	Source   string
	NewsList []entity.News
}

const pathToTemplate = "cli/internal/template/news.tmpl"

// Apply the template to the news and prints the results.

func (t Template) CreateTemplate(keywords string) *template.Template {
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

	tmpl, err := template.New("news").Funcs(funcMap).ParseFiles(pathToTemplate) // Update with the actual path
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return tmpl
}

// Prepare the template data for rendering.
func (t Template) Prepare() Template {
	groupedMap := group(t.News)
	var groupedList []*groupedNews
	sourceList := make([]string, 0)
	for el := groupedMap.Oldest(); el != nil; el = el.Next() {
		source := el.Key.(string)
		sourceList = append(sourceList, source)
		newsList := el.Value.([]entity.News)
		groupedList = append(groupedList, &groupedNews{Source: source, NewsList: newsList})
	}
	return Template{
		Filters:   "sources=" + strings.Join(sourceList, ",") + " ",
		News:      t.News,
		Criterion: t.Criterion,
		Grouped:   groupedList,
	}
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
