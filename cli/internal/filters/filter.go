package filters

import (
	"news-aggregator/internal/entity"
	"news-aggregator/internal/validator"
	"strings"
	"time"
)

// NewsFilter is a filtering of news according to specified parameters.
type NewsFilter interface {
	Filter(news []entity.News) []entity.News
	String() string
}

// InitializeFilters based on provided parameters.
func InitializeFilters(keywords, dateStart, dateEnd *string) []NewsFilter {
	var newsFilters []NewsFilter

	if keywordFilter := convertKeywords(keywords); keywordFilter != nil {
		newsFilters = append(newsFilters, keywordFilter)
	}
	if dateStartFilter := convertDateStart(dateStart); dateStartFilter != nil {
		newsFilters = append(newsFilters, dateStartFilter)
	}
	if dateEndFilter := convertDateEnd(dateEnd); dateEndFilter != nil {
		newsFilters = append(newsFilters, dateEndFilter)
	}

	return newsFilters
}

// convertKeywords a comma-separated keyword string into a filter.
func convertKeywords(keywords *string) *Keyword {
	if len(*keywords) > 0 {
		keywordList := strings.Split(*keywords, ",")
		return &Keyword{Keywords: keywordList}
	}
	return nil
}

// convertDateStart string into a DateStart filter.
func convertDateStart(dateStart *string) *DateStart {
	if len(*dateStart) > 0 {
		startDate, _ := time.Parse(validator.DateFormat, *dateStart)
		return &DateStart{StartDate: startDate}
	}
	return nil
}

// convertDateEnd string into a DateEnd filter.
func convertDateEnd(dateEnd *string) *DateEnd {
	if len(*dateEnd) > 0 {
		endDate, _ := time.Parse(validator.DateFormat, *dateEnd)
		return &DateEnd{EndDate: endDate}
	}
	return nil
}
