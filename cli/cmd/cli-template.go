package main

import (
	"NewsAggregator/cli/internal/entity"
	"fmt"
	"github.com/wk8/go-ordered-map"
	"log"
	"os"
	"strings"
	"text/template"
)

// Template represents the data structure for the news template.
type Template struct {
	Filters   string
	Count     int
	News      []entity.News
	Keywords  string
	Criterion string
	Order     string
	Grouped   []*groupedNews
}

// groupedNews represents a group of news items.
type groupedNews struct {
	Source   string
	NewsList []entity.News
}

// apply the template to the news and prints the results.
func (t Template) apply(sourceList []string) {

	funcMap := template.FuncMap{
		"highlight": func(text, keywords string) string {
			if keywords != "" {
				for _, keyword := range strings.Split(keywords, ",") {
					text = strings.ReplaceAll(text, keyword, "~~"+keyword+"~~")
				}
			}
			return text
		},
		"toString": func(v interface{}) string {
			return fmt.Sprintf("%v", v)
		},
	}

	tmpl, err := template.New("news").Funcs(funcMap).ParseFiles("cli/internal/entity/news.tmpl")
	if err != nil {
		log.Fatal(err)
		return
	}

	fil := fmt.Sprintf("sources=%s, keywords=%s, sort-by=%s, order=%s", strings.Join(sourceList, ","), t.Keywords, t.Criterion, t.Order)
	data := t.prepare(fil)

	err = tmpl.ExecuteTemplate(os.Stdout, "news", data)
	if err != nil {
		panic(err)
	}
}

// prepare  the template data for rendering.
func (t Template) prepare(filters string) Template {
	groupedMap := group(t.News)
	var groupedList []*groupedNews

	for el := groupedMap.Oldest(); el != nil; el = el.Next() {
		source := el.Key.(string)
		newsList := el.Value.([]entity.News)
		groupedList = append(groupedList, &groupedNews{Source: source, NewsList: newsList})
	}

	return Template{
		Filters:   filters,
		Count:     len(t.News),
		News:      t.News,
		Keywords:  t.Keywords,
		Criterion: t.Criterion,
		Order:     t.Order,
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
