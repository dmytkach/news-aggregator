package internal

import (
	"NewsAggregator/cli/internal/entity"
	"sort"
)

// Sort news according to the specified criteria and order.
func Sort(news []entity.News, criterion, order string) []entity.News {
	sort.Slice(news, func(i, j int) bool {
		if criterion == "date" {
			if order == "DESC" {
				return news[i].Date.After(news[j].Date)
			}
			return news[i].Date.Before(news[j].Date)
		} else if criterion == "source" {
			if order == "DESC" {
				return news[i].Source > news[j].Source
			}
			return news[i].Source < news[j].Source
		}
		return false
	})
	return news
}
