package initializers

import (
	"news-aggregator/internal/filters"
	"news-aggregator/internal/validator"
	"strings"
	"time"
)

// InitializeFilters based on provided parameters.
func InitializeFilters(keywords, dateStart, dateEnd *string) []filters.NewsFilter {
	var newsFilters []filters.NewsFilter

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
func convertKeywords(keywords *string) *filters.Keyword {
	if len(*keywords) > 0 {
		keywordList := strings.Split(*keywords, ",")
		return &filters.Keyword{Keywords: keywordList}
	}
	return nil
}

// convertDateStart string into a DateStart filter.
func convertDateStart(dateStart *string) *filters.DateStart {
	if len(*dateStart) > 0 {
		startDate, _ := time.Parse(validator.DateFormat, *dateStart)
		return &filters.DateStart{StartDate: startDate}
	}
	return nil
}

// convertDateEnd string into a DateEnd filter.
func convertDateEnd(dateEnd *string) *filters.DateEnd {
	if len(*dateEnd) > 0 {
		endDate, _ := time.Parse(validator.DateFormat, *dateEnd)
		return &filters.DateEnd{EndDate: endDate}
	}
	return nil
}
