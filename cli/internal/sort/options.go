package sort

import (
	"news-aggregator/internal/entity"
	"sort"
)

const Desc = "DESC"

type Options struct {
	Criterion string
	Order     string
}

// Apply news according to the specified criteria and order.
func (s *Options) Apply(news []entity.News) []entity.News {
	sort.Slice(news, func(i, j int) bool {
		if s.Criterion == "date" {
			if s.Order == Desc {
				return news[i].Date.After(news[j].Date)
			}
			return news[i].Date.Before(news[j].Date)
		} else if s.Criterion == "source" {
			if s.Order == Desc {
				return news[i].Source > news[j].Source
			}
			return news[i].Source < news[j].Source
		}
		return false
	})
	return news
}
