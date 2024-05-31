package filters

import (
	"NewsAggregator/cli/internal/entity"
	"time"
)

// DateStart filters news from the specified date.
type DateStart struct {
	StartDate time.Time
}

// Filter filters news starting from the specified date.
func (dsf *DateStart) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		if item.Date.After(dsf.StartDate) || item.Date.Equal(dsf.StartDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
