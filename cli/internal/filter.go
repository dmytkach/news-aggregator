package internal

import (
	"news-aggregator/cli/internal/entity"
)

// NewsFilter is a filtering of news according to specified parameters.
type NewsFilter interface {
	Filter(news []entity.News) []entity.News
	String() string
}
