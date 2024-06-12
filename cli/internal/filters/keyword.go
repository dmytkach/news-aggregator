package filters

import (
	"NewsAggregator/cli/internal/entity"
	"strings"
)

// Keyword filters news by keywords.
type Keyword struct {
	Keywords []string
}

// Filter news by keywords in the title and description.
func (k *Keyword) Filter(news []entity.News) []entity.News {
	var filtered []entity.News
	for _, item := range news {
		for _, keyword := range k.Keywords {
			if strings.Contains(strings.ToLower(string(item.Title)), strings.ToLower(keyword)) ||
				strings.Contains(strings.ToLower(string(item.Description)), strings.ToLower(keyword)) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}
func (k *Keyword) String() string {
	return "keywords=" + strings.Join(k.Keywords, ",")
}
