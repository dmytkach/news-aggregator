package handlers

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mocks"
	"testing"
)

func TestSourcesGet(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	expectedSources := []entity.Source{{Name: "bbc_news", PathsToFile: nil}}
	mockSourceManager.On("GetSources").Return(expectedSources, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `[{"Name":"bbc_news","PathsToFile":null}]`, rr.Body.String(), "Response body does not match expected")
}
func TestGetSourcesEmpty(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.On("GetSources").Return([]entity.Source{}, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `[]`, rr.Body.String(), "Response body does not match expected")
}

func TestDownloadSource(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	feed := entity.Feed{Name: "test_feed"}
	mockFeedManager.On("FetchFeed", "http://example.com/feed").Return(feed, nil)
	expectedSource := entity.Source{Name: "test_feed", PathsToFile: []entity.PathToFile{"http://example.com/feed"}}
	mockSourceManager.On("CreateSource", "test_feed", "http://example.com/feed").Return(expectedSource, nil)

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
	assert.JSONEq(t, `{"Name":"test_feed","PathsToFile":["http://example.com/feed"]}`, rr.Body.String(), "Response body does not match expected")
}
func TestDownloadSourceError(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockFeedManager.On("FetchFeed", "http://example.com/feed").Return(entity.Feed{}, errors.New("fetch error"))

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestDownloadSourceMissingURL(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

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
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	feed := entity.Feed{Name: "test_feed"}
	mockFeedManager.On("FetchFeed", "http://example.com/feed").Return(feed, nil)
	mockSourceManager.On("CreateSource", "test_feed", "http://example.com/feed").Return(entity.Source{}, errors.New("creation error"))

	req, err := http.NewRequest("POST", "/sources?url=http://example.com/feed", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}
func TestUpdateSource(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.On("UpdateSource", "http://oldurl.com", "http://newurl.com").Return(nil)

	req, err := http.NewRequest("PUT", "/sources?oldUrl=http://oldurl.com&newUrl=http://newurl.com", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}
func TestRemoveSource(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.On("RemoveSourceByName", "bbc_news").Return(nil)

	req, err := http.NewRequest("DELETE", "/sources?name=bbc_news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")
}
func TestUpdateSourceError(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.On("UpdateSource", "http://oldurl.com", "http://newurl.com").Return(errors.New("update error"))

	req, err := http.NewRequest("PUT", "/sources?oldUrl=http://oldurl.com&newUrl=http://newurl.com", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestRemoveSourceMissingName(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

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
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

	sourceHandler := SourceHandler{
		SourceManager: mockSourceManager,
		FeedManager:   mockFeedManager,
	}

	mockSourceManager.On("RemoveSourceByName", "bbc_news").Return(errors.New("remove error"))

	req, err := http.NewRequest("DELETE", "/sources?name=bbc_news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestSourcesMethodNotAllowed(t *testing.T) {
	mockSourceManager := new(mocks.MockSourceManager)
	mockFeedManager := new(mocks.MockFeedManager)

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
