package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mock_managers"
)

func TestSourcesGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	expectedSources := []entity.Source{{Name: "bbc_news", PathsToFile: nil}}
	mockSourceManager.EXPECT().GetSources().Return(expectedSources, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `[{"Name":"bbc_news","PathsToFile":null}]`, rr.Body.String(), "Response body does not match expected")
}

func TestGetSourcesEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.EXPECT().GetSources().Return([]entity.Source{}, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `[]`, rr.Body.String(), "Response body does not match expected")
}

func TestDownloadSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	feed := entity.Feed{Name: "test_feed"}
	mockFeedManager.EXPECT().FetchFeed("http://example.com/feed").Return(feed, nil)
	expectedSource := entity.Source{Name: "test_feed", PathsToFile: []entity.PathToFile{"http://example.com/feed"}}
	mockSourceManager.EXPECT().CreateSource("test_feed", "http://example.com/feed").Return(expectedSource, nil)

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `{"Name":"test_feed","PathsToFile":["http://example.com/feed"]}`, rr.Body.String(), "Response body does not match expected")
}

func TestDownloadSourceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockFeedManager.EXPECT().FetchFeed("http://example.com/feed").Return(entity.Feed{}, errors.New("fetch error"))

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestDownloadSourceMissingURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	req, err := http.NewRequest("POST", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}

func TestDownloadSourceCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	feed := entity.Feed{Name: "test_feed"}
	mockFeedManager.EXPECT().FetchFeed("http://example.com/feed").Return(feed, nil)
	mockSourceManager.EXPECT().CreateSource("test_feed", "http://example.com/feed").Return(entity.Source{}, errors.New("creation error"))

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestUpdateSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.EXPECT().UpdateSource("http://oldurl.com", "http://newurl.com").Return(nil)

	req, err := http.NewRequest("PUT", "/sources?oldUrl=http://oldurl.com&newUrl=http://newurl.com", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}

func TestRemoveSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.EXPECT().RemoveSourceByName("bbc_news").Return(nil)

	req, err := http.NewRequest("DELETE", "/sources?name=bbc_news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}

func TestUpdateSourceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.EXPECT().UpdateSource("http://oldurl.com", "http://newurl.com").Return(errors.New("update error"))

	req, err := http.NewRequest("PUT", "/sources?oldUrl=http://oldurl.com&newUrl=http://newurl.com", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestRemoveSourceMissingName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	req, err := http.NewRequest("DELETE", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}

func TestRemoveSourceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.EXPECT().RemoveSourceByName("bbc_news").Return(errors.New("remove error"))

	req, err := http.NewRequest("DELETE", "/sources?name=bbc_news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestSourcesMethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	mockFeedManager := mock_managers.NewMockFeedManager(ctrl)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	req, err := http.NewRequest("PATCH", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected status Method Not Allowed")
}
