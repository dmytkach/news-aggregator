package internal

import (
	"NewsAggregator/cli/internal/entity"
	"sort"
	"strings"
)

func SortNews(news []entity.News, sortBy, sortOrder string) {
	switch strings.ToLower(sortBy) {
	case "source":
		sort.Slice(news, func(i, j int) bool {
			if strings.ToLower(sortOrder) == "asc" {
				return news[i].Link < news[j].Link
			}
			return news[i].Link > news[j].Link
		})
	case "date":
		fallthrough
	default:
		sort.Slice(news, func(i, j int) bool {
			if strings.ToLower(sortOrder) == "asc" {
				return news[i].Date.Before(news[j].Date)
			}
			return news[i].Date.After(news[j].Date)
		})
	}
}
