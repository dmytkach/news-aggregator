package service

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mock_managers"
)

func TestFetchService_UpdateNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockNewsManager := mock_managers.NewMockNewsManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sources := []entity.Source{
		{
			Name: "Source1",
			PathsToFile: []entity.PathToFile{
				"file1.xml",
				"file2.xml",
			},
		},
	}

	mockSourceManager.EXPECT().GetSources().Return(sources, nil).Times(1)
	mockNewsManager.EXPECT().GetNewsFromFolder("Source1").Return([]entity.News{
		{Link: "link1"},
	}, nil).Times(2)
	mockFeedManager.EXPECT().FetchFeed("file1.xml").Return(entity.Feed{
		Name: "Feed1",
		News: []entity.News{
			{Link: "link2"},
		},
	}, nil).Times(1)
	mockFeedManager.EXPECT().FetchFeed("file2.xml").Return(entity.Feed{
		Name: "Feed2",
		News: []entity.News{
			{Link: "link3"},
		},
	}, nil).Times(1)
	mockNewsManager.EXPECT().AddNews([]entity.News{
		{Link: "link2"},
	}, "Feed1").Return(nil).Times(1)
	mockNewsManager.EXPECT().AddNews([]entity.News{
		{Link: "link3"},
	}, "Feed2").Return(nil).Times(1)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   mockFeedManager,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err)
}

func TestFetchService_UpdateNews_FetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockNewsManager := mock_managers.NewMockNewsManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sources := []entity.Source{
		{
			Name: "Source1",
			PathsToFile: []entity.PathToFile{
				"file1.xml",
			},
		},
	}

	mockSourceManager.EXPECT().GetSources().Return(sources, nil).Times(1)
	mockFeedManager.EXPECT().FetchFeed("file1.xml").Return(entity.Feed{}, errors.New("fetch error")).Times(1)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   mockFeedManager,
	}

	err := fetchService.UpdateNews()
	assert.NoError(t, err, "Expected no error from UpdateNews")
}

func TestFetchService_fetchNewsFromSource_FetchError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockNewsManager := mock_managers.NewMockNewsManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	resource := entity.Source{
		Name: "Source1",
		PathsToFile: []entity.PathToFile{
			"file1.xml",
		},
	}

	mockFeedManager.EXPECT().FetchFeed("file1.xml").Return(entity.Feed{}, errors.New("fetch error")).Times(1)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   mockFeedManager,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource due to continue on fetch error")
}

func TestFetchService_fetchNewsFromSource_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockNewsManager := mock_managers.NewMockNewsManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

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

	mockFeedManager.EXPECT().FetchFeed("file1.xml").Return(newFeed, nil).Times(1)
	mockNewsManager.EXPECT().GetNewsFromFolder("Source1").Return(existingNews, nil).Times(1)
	mockNewsManager.EXPECT().AddNews([]entity.News{{Link: "new_link"}}, "Source1").Return(nil).Times(1)

	fetchService := Fetch{
		SourceManager: mockSourceManager,
		NewsManager:   mockNewsManager,
		FeedManager:   mockFeedManager,
	}

	err := fetchService.fetchNewsFromSource(resource)
	assert.NoError(t, err, "Expected no error from fetchNewsFromSource")
}
