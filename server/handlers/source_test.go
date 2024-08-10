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
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	expectedSources := []entity.Source{{Name: "bbc_news", PathToFile: "test-path-to-file"}}
	mockSourceManager.EXPECT().GetSources().Return(expectedSources, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[{"Name":"bbc_news","PathToFile":"test-path-to-file"}]`, rr.Body.String())
}

func TestGetSourceByName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	expectedSource := entity.Source{Name: "bbc_news", PathToFile: "test-path-to-file"}
	mockSourceManager.EXPECT().GetSource("bbc_news").Return(expectedSource, nil)

	req, err := http.NewRequest("GET", "/sources?name=bbc_news", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"Name":"bbc_news","PathToFile":"test-path-to-file"}`, rr.Body.String())
}

func TestGetSourcesEmpty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().GetSources().Return([]entity.Source{}, nil)

	req, err := http.NewRequest("GET", "/sources", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[]`, rr.Body.String())
}

func TestDownloadSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	expectedSource := entity.Source{Name: "test_feed", PathToFile: "http://example.com/feed"}
	mockSourceManager.EXPECT().CreateSource("test_feed", "http://example.com/feed").Return(expectedSource, nil)

	req, err := http.NewRequest("POST", "/sources?name=test_feed&url=http://example.com/feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `{"Name":"test_feed","PathToFile":"http://example.com/feed"}`, rr.Body.String())
}

func TestDownloadSourceMissingParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	req, err := http.NewRequest("POST", "/sources", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDownloadSourceCreationError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().CreateSource("test_feed", "http://example.com/feed").Return(entity.Source{}, errors.New("creation error"))

	req, err := http.NewRequest("POST", "/sources?name=test_feed&url=http://example.com/feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestUpdateSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().UpdateSource("test_feed", "http://newurl.com").Return(nil)

	req, err := http.NewRequest("PUT", "/sources?name=test_feed&newUrl=http://newurl.com", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUpdateSourceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().UpdateSource("test_feed", "http://newurl.com").Return(errors.New("update error"))

	req, err := http.NewRequest("PUT", "/sources?name=test_feed&newUrl=http://newurl.com", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestRemoveSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().RemoveSourceByName("test_feed").Return(nil)

	req, err := http.NewRequest("DELETE", "/sources?name=test_feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRemoveSourceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	mockSourceManager.EXPECT().RemoveSourceByName("test_feed").Return(errors.New("remove error"))

	req, err := http.NewRequest("DELETE", "/sources?name=test_feed", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

func TestSourcesMethodNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	sourceHandler := SourceHandler{SourceManager: mockSourceManager}

	req, err := http.NewRequest("PATCH", "/sources", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(sourceHandler.Sources)
	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code)
}
