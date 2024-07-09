package template

import (
	"fmt"
	"github.com/wk8/go-ordered-map"
	"log"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/sort"
	"regexp"
	"strings"
	"text/template"
)

// Data represents the data structure for the news template.
type Data struct {
	Header  Header
	News    []entity.News
	Grouped []*groupedNews
}
type Header struct {
	Sources     string
	Filters     string
	SortOptions sort.Options
}

type groupedNews struct {
	Source   string
	NewsList []entity.News
}

const pathToTemplate = "cli/internal/template/news.tmpl"

// Create generates a template for rendering news.
func (t Data) Create(keywords string) (*template.Template, error) {
	funcMap := template.FuncMap{
		"highlight": func(text string) string {
			if len(keywords) == 0 {
				return text
			}
			for _, keyword := range strings.Split(keywords, ",") {
				re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(keyword))
				text = re.ReplaceAllStringFunc(text, func(matched string) string {
					return "~~" + matched + "~~"
				})
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
func (t Data) Prepare() Data {
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

// group categorizes news items by their source.
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
