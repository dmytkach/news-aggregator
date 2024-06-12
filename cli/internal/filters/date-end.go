package filters

import (
	"NewsAggregator/cli/internal/entity"
	"time"
)

// DateEnd filters news up to the specified date.
type DateEnd struct {
	EndDate time.Time
}

// Filter news up to a specified date.
func (de *DateEnd) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		if item.Date.Before(de.EndDate) || item.Date.Equal(de.EndDate) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}
func (de *DateEnd) String() string {
	return "date-end=" + de.EndDate.Format("2006-01-02")
}
