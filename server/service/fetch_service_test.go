package service

import (
	"errors"
	"news-aggregator/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSourceManager struct {
	mock.Mock
}

func (m *MockSourceManager) CreateSource(name, url string) (entity.Source, error) {
	panic("should not be called")
}

func (m *MockSourceManager) GetSource(name string) (entity.Source, error) {
	panic("should not be called")
}

func (m *MockSourceManager) UpdateSource(oldUrl, newUrl string) error {
	panic("should not be called")
}

func (m *MockSourceManager) RemoveSourceByName(sourceName string) error {
	panic("should not be called")
}

func (m *MockSourceManager) GetSources() ([]entity.Source, error) {
	args := m.Called()
	return args.Get(0).([]entity.Source), args.Error(1)
}

type MockNewsManager struct {
	mock.Mock
}

func (m *MockNewsManager) GetNewsSourceFilePath(sourceName []string) (map[string][]string, error) {
	panic("should not be called")
}

func (m *MockNewsManager) GetNewsFromFolder(folderName string) ([]entity.News, error) {
	args := m.Called(folderName)
	return args.Get(0).([]entity.News), args.Error(1)
}

func (m *MockNewsManager) AddNews(news []entity.News, folderName string) error {
	args := m.Called(news, folderName)
	return args.Error(0)
}

type MockFetch struct {
	mock.Mock
}

func (m *MockFetch) Fetch(path string) (entity.Feed, error) {
	args := m.Called(path)
	return args.Get(0).(entity.Feed), args.Error(1)
}

func TestFetchService_UpdateNews(t *testing.T) {
	mockSourceManager := new(MockSourceManager)
	mockNewsManager := new(MockNewsManager)
	mockFetch := new(MockFetch)

	sources := []entity.Source{
		{
			Name: "Source1",
			PathsToFile: []entity.PathToFile{
				"file1.xml",
				"file2.xml",
			},
		},
	}

	mockSourceManager.On("GetSources").Return(sources, nil)
	mockNewsManager.On("GetNewsFromFolder", "Source1").Return([]entity.News{}, nil)
	mockFetch.On("Fetch", "file1.xml").Return(entity.Feed{
		Name: "Feed1",
		News: []entity.News{
			{Link: "link1"},
			{Link: "link2"},
		},
	}, nil)
	mockFetch.On("Fetch", "file2.xml").Return(entity.Feed{
		Name: "Feed2",
		News: []entity.News{
			{Link: "link3"},
			{Link: "link4"},
		},
	}, nil)
	mockNewsManager.On("AddNews", mock.Anything, "Feed1").Return(nil)
	mockNewsManager.On("AddNews", mock.Anything, "Feed2").Return(nil)

	fetchService := FetchService{
		SourceRepo: mockSourceManager,
		NewsRepo:   mockNewsManager,
		Fetch:      mockFetch,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err)

	mockSourceManager.AssertExpectations(t)
	mockNewsManager.AssertExpectations(t)
	mockFetch.AssertExpectations(t)
}

func TestFetchService_UpdateNews_FetchError(t *testing.T) {
	mockSourceManager := new(MockSourceManager)
	mockNewsManager := new(MockNewsManager)
	mockFetch := new(MockFetch)

	sources := []entity.Source{
		{
			Name: "Source1",
			PathsToFile: []entity.PathToFile{
				"file1.xml",
			},
		},
	}

	mockSourceManager.On("GetSources").Return(sources, nil)
	mockFetch.On("Fetch", "file1.xml").Return(entity.Feed{}, errors.New("fetch error"))

	fetchService := FetchService{
		SourceRepo: mockSourceManager,
		NewsRepo:   mockNewsManager,
		Fetch:      mockFetch,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err, "Expected no error from UpdateNews")

	mockSourceManager.AssertExpectations(t)
	mockFetch.AssertExpectations(t)
}
func TestFetchService_fetchNewsFromSource_FetchError(t *testing.T) {
	mockSourceManager := new(MockSourceManager)
	mockNewsManager := new(MockNewsManager)
	mockFetch := new(MockFetch)

	resource := entity.Source{
		Name: "Source1",
		PathsToFile: []entity.PathToFile{
			"file1.xml",
		},
	}

	mockFetch.On("Fetch", "file1.xml").Return(entity.Feed{}, errors.New("fetch error"))

	fetchService := FetchService{
		SourceRepo: mockSourceManager,
		NewsRepo:   mockNewsManager,
		Fetch:      mockFetch,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource due to continue on fetch error")

	mockFetch.AssertExpectations(t)
}

func TestFetchService_fetchNewsFromSource_Success(t *testing.T) {
	mockSourceManager := new(MockSourceManager)
	mockNewsManager := new(MockNewsManager)
	mockFetch := new(MockFetch)

	resource := entity.Source{
		Name: "Source1",
		PathsToFile: []entity.PathToFile{
			"file1.xml",
		},
	}

	existingNews := []entity.News{
		{
			Link: "existing_link",
		},
	}

	newFeed := entity.Feed{
		Name: "Source1",
		News: []entity.News{
			{
				Link: "new_link",
			},
		},
	}

	mockFetch.On("Fetch", "file1.xml").Return(newFeed, nil)
	mockNewsManager.On("GetNewsFromFolder", "Source1").Return(existingNews, nil)
	mockNewsManager.On("AddNews", []entity.News{{Link: "new_link"}}, "Source1").Return(nil)

	fetchService := FetchService{
		SourceRepo: mockSourceManager,
		NewsRepo:   mockNewsManager,
		Fetch:      mockFetch,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource")

	mockFetch.AssertExpectations(t)
	mockNewsManager.AssertExpectations(t)
}
