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
func (ds *DateStart) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		if item.Date.After(ds.StartDate) || item.Date.Equal(ds.StartDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
func (ds *DateStart) String() string {
	return " date-start=" + ds.StartDate.Format("2006-01-02")
}
