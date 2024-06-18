package internal

import (
	"news-aggregator/internal/entity"
)

type Sort interface {
	Apply(news []entity.News) []entity.News
}
