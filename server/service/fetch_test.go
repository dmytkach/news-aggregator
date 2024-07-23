package service

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mocks"
	"testing"
)

func TestFetchService_UpdateNews(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockNewsManager := new(mocks.MockNewsManager)
	mockFeedManager := new(mocks.MockFeedManager)

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
	mockFeedManager.On("FetchFeed", "file1.xml").Return(entity.Feed{
		Name: "Feed1",
		News: []entity.News{
			{Link: "link1"},
			{Link: "link2"},
		},
	}, nil)
	mockFeedManager.On("FetchFeed", "file2.xml").Return(entity.Feed{
		Name: "Feed2",
		News: []entity.News{
			{Link: "link3"},
			{Link: "link4"},
		},
	}, nil)
	mockNewsManager.On("AddNews", mock.Anything, "Feed1").Return(nil)
	mockNewsManager.On("AddNews", mock.Anything, "Feed2").Return(nil)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   mockFeedManager,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err)

	mockSourceManager.AssertExpectations(t)
	mockNewsManager.AssertExpectations(t)
	mockFeedManager.AssertExpectations(t)
}

func TestFetchService_UpdateNews_FetchError(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockNewsManager := new(mocks.MockNewsManager)
	MockFeedManager := new(mocks.MockFeedManager)

	sources := []entity.Source{
		{
			Name: "Source1",
			PathsToFile: []entity.PathToFile{
				"file1.xml",
			},
		},
	}

	mockSourceManager.On("GetSources").Return(sources, nil)
	MockFeedManager.On("FetchFeed", "file1.xml").Return(entity.Feed{}, errors.New("fetch error"))

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   MockFeedManager,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err, "Expected no error from UpdateNews")

	mockSourceManager.AssertExpectations(t)
	MockFeedManager.AssertExpectations(t)
}
func TestFetchService_fetchNewsFromSource_FetchError(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockNewsManager := new(mocks.MockNewsManager)
	MockFeedManager := new(mocks.MockFeedManager)

	resource := entity.Source{
		Name: "Source1",
		PathsToFile: []entity.PathToFile{
			"file1.xml",
		},
	}

	MockFeedManager.On("FetchFeed", "file1.xml").Return(entity.Feed{}, errors.New("fetch error"))

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   MockFeedManager,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource due to continue on fetch error")

	MockFeedManager.AssertExpectations(t)
}

func TestFetchService_fetchNewsFromSource_Success(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockNewsManager := new(mocks.MockNewsManager)
	MockFeedManager := new(mocks.MockFeedManager)

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

	MockFeedManager.On("FetchFeed", "file1.xml").Return(newFeed, nil)
	mockNewsManager.On("GetNewsFromFolder", "Source1").Return(existingNews, nil)
	mockNewsManager.On("AddNews", []entity.News{{Link: "new_link"}}, "Source1").Return(nil)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   MockFeedManager,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource")

	MockFeedManager.AssertExpectations(t)
	mockNewsManager.AssertExpectations(t)
}
