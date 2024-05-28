package filter

import (
	"NewsAggregator/cli/internal/entity"
	"strings"
)

// KeywordFilter filters news by keywords.
type KeywordFilter struct {
	Keywords []string
}

// Filter filters news by keywords in the title and description.
func (kf *KeywordFilter) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		for _, keyword := range kf.Keywords {
			if strings.Contains(strings.ToLower(string(item.Title)), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(string(item.Description)), strings.ToLower(keyword)) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}
