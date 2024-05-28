package filter

import (
	"NewsAggregator/cli/internal/entity"
	"time"
)

// DateEndFilter filters news up to the specified date.
type DateEndFilter struct {
	EndDate time.Time
}

// Filter filters news up to a specified date.
func (def *DateEndFilter) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		if item.Date.Before(def.EndDate) || item.Date.Equal(def.EndDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
