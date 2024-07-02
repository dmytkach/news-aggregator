package internal

import (
	"news-aggregator/internal/entity"
)

// Sort provides the ability to sort news articles.
type Sort interface {
	Sort(news []entity.News) []entity.News
}
