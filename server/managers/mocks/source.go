package mocks

import (
	"github.com/stretchr/testify/mock"
	"news-aggregator/internal/entity"
)

type MockSourceManager struct {
	mock.Mock
}

func (m *MockSourceManager) CreateSource(name, url string) (entity.Source, error) {
	args := m.Called(name, url)
	return args.Get(0).(entity.Source), args.Error(1)
}

func (m *MockSourceManager) GetSource(name string) (entity.Source, error) {
	args := m.Called(name)
	return args.Get(0).(entity.Source), args.Error(1)
}

func (m *MockSourceManager) GetSources() ([]entity.Source, error) {
	args := m.Called()
	return args.Get(0).([]entity.Source), args.Error(1)
}

func (m *MockSourceManager) UpdateSource(oldUrl, newUrl string) error {
	args := m.Called(oldUrl, newUrl)
	return args.Error(0)
}

func (m *MockSourceManager) RemoveSourceByName(sourceName string) error {
	args := m.Called(sourceName)
	return args.Error(0)
}
