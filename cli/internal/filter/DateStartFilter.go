package filter

import (
	"NewsAggregator/cli/internal/entity"
	"time"
)

type DateStartFilter struct {
	StartDate time.Time
}

// Filter filters news starting from the specified date.
func (dsf *DateStartFilter) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		if item.Date.After(dsf.StartDate) || item.Date.Equal(dsf.StartDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
