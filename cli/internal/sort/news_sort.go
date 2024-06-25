package sort

import (
	"news-aggregator/internal/entity"
	"sort"
	"strings"
)

// desc sort order
const desc = "DESC"
const asc = "ASC"

// Options for the sorting process
type Options struct {
	Criterion string
	Order     string
}

// Sort news according to the specified Options.
func (s *Options) Sort(news []entity.News) []entity.News {
	if s.Criterion == "" {
		s.Criterion = "date"
	}
	sort.Slice(news, func(i, j int) bool {
		if strings.EqualFold(s.Criterion, "date") {
			if s.Order == desc {
				return news[i].Date.After(news[j].Date)
			} else if s.Order == asc {
				return news[i].Date.Before(news[j].Date)
			}
		} else if strings.EqualFold(s.Criterion, "source") {
			if s.Order == desc {
				return news[i].Source > news[j].Source
			} else if s.Order == asc {
				return news[i].Source < news[j].Source
			}
		}
		return false
	})
	return news
}
