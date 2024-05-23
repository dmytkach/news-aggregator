package filter

import (
	"NewsAggregator/cli/internal/entity"
	"strings"
	"time"
)

// NewsFilter is a filtering of news according to specified parameters.
type NewsFilter struct {
	News *[]entity.News
}

// FilterByKeywords filters news items based on the provided keywords in a title.
func (newsFilter *NewsFilter) FilterByKeywords(keywords []string) []entity.News {
	var filtered []entity.News
	for _, news := range *newsFilter.News {
		for _, keyword := range keywords {
			if strings.Contains(strings.ToLower(string(news.Title)), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(string(news.Description)), strings.ToLower(keyword)) {
				filtered = append(filtered, news)
				break
			}
		}
	}
	return filtered
}

// FilterByDate filters news items within the specified date range.
func (newsFilter *NewsFilter) FilterByDate(startDate, endDate time.Time) []entity.News {
	var filtered []entity.News
	for _, news := range *newsFilter.News {
		if news.Date.After(startDate) && news.Date.Before(endDate) {
			filtered = append(filtered, news)
		}
	}
	return filtered
}
