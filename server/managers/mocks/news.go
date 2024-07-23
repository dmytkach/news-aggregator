package mocks

import (
	"github.com/stretchr/testify/mock"
	"news-aggregator/internal/entity"
)

type MockNewsManager struct {
	mock.Mock
}

func (m *MockNewsManager) AddNews(newsToAdd []entity.News, newsSource string) error {
	args := m.Called(newsToAdd, newsSource)
	return args.Error(0)
}

func (m *MockNewsManager) GetNewsFromFolder(folderName string) ([]entity.News, error) {
	args := m.Called(folderName)
	return args.Get(0).([]entity.News), args.Error(1)
}

func (m *MockNewsManager) GetNewsSourceFilePath(sourceNames []string) (map[string][]string, error) {
	args := m.Called(sourceNames)
	return args.Get(0).(map[string][]string), args.Error(1)
}
