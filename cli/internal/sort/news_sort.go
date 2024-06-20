package sort

import (
	"news-aggregator/internal/entity"
	"sort"
)

// desc sort order
const desc = "DESC"

// Options for the sorting process
type Options struct {
	Criterion string
	Order     string
}

// Sort news according to the specified Options.
func (s *Options) Sort(news []entity.News) []entity.News {
	sort.Slice(news, func(i, j int) bool {
		if s.Criterion == "date" {
			if s.Order == desc {
				return news[i].Date.After(news[j].Date)
			}
			return news[i].Date.Before(news[j].Date)
		} else if s.Criterion == "source" {
			if s.Order == desc {
				return news[i].Source > news[j].Source
			}
			return news[i].Source < news[j].Source
		}
		return false
	})
	return news
}
