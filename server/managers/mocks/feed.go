package mocks

import (
	"github.com/stretchr/testify/mock"
	"news-aggregator/internal/entity"
)

type MockFeedManager struct {
	mock.Mock
}

func (m *MockFeedManager) FetchFeed(path string) (entity.Feed, error) {
	args := m.Called(path)
	return args.Get(0).(entity.Feed), args.Error(1)
}
