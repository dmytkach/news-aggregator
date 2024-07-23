package handlers

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mocks"
	"testing"
	"time"
)

func setupNewsHandlerTest() (*NewsHandler, *mocks.MockNewsManager, *mocks.MockSourceManager) {
	mockNewsManager := new(mocks.MockNewsManager)
	mockSourceManager := new(mocks.MockSourceManager)
	handler := &NewsHandler{
		NewsManager:   mockNewsManager,
		SourceManager: mockSourceManager,
	}
	return handler, mockNewsManager, mockSourceManager
}
func TestNewsHandler(t *testing.T) {
	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest()

	mockSources := []entity.Source{
		{Name: "bbc_news"},
	}
	mockSourceManager.On("GetSources").Return(mockSources, nil)

	mockNewsManager.On("GetNewsSourceFilePath", []string{"bbc_news"}).Return(map[string][]string{
		"bbc_news": {"../../internal/testdata/bbc_news/ready_news.json"},
	}, nil)

	req, err := http.NewRequest("GET", "/news?sources=bbc_news&keywords=England", nil)
	assert.Nil(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected status OK")

	expected := []entity.News{
		{
			Title:       "Watch England fans go wild as Bellingham scores late equaliser",
			Description: "Watch England fans erupt at a fanpark in Wembley as Jude Bellingham scores a late stunner to send the game to extra time against Slovakia in the Euro 2024 last-16 match in Gelsenkirchen.\n\n\n",
			Link:        "https://www.bbc.com/sport/football/videos/cl4yj1ve5z7o",
			Source:      "bbc_news",
			Date:        time.Date(2024, 6, 30, 19, 31, 26, 0, time.UTC),
		},
	}
	var actual []entity.News
	err = json.NewDecoder(rr.Body).Decode(&actual)
	assert.Nil(t, err, "Expected no error decoding response body")
	assert.ElementsMatch(t, expected, actual, "Expected response body to match")
}
func TestNewsHandlerInvalidSource(t *testing.T) {
	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest()

	var mockSources []entity.Source
	mockSourceManager.On("GetSources").Return(mockSources, nil)

	mockNewsManager.On("GetNewsSourceFilePath", mock.Anything).Return(map[string][]string{}, errors.New("error getting news source file paths"))
	handler.NewsManager = mockNewsManager

	req, err := http.NewRequest("GET", "/news?sources=invalid_source", nil)
	assert.Nil(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}
func TestNewsHandlerErrorGettingSources(t *testing.T) {
	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest()

	mockSourceManager.On("GetSources").Return([]entity.Source{}, errors.New("error getting sources"))

	mockNewsManager.On("GetNewsSourceFilePath", mock.Anything).Return(map[string][]string{}, nil)

	req, err := http.NewRequest("GET", "/news?sources=invalid_source", nil)
	assert.Nil(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestNewsHandlerErrorGettingNewsSourceFilePath(t *testing.T) {
	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest()

	mockSources := []entity.Source{
		{Name: "bbc_news"},
	}
	mockSourceManager.On("GetSources").Return(mockSources, nil)
	mockNewsManager.On("GetNewsSourceFilePath", []string{"bbc_news"}).Return(map[string][]string{}, errors.New("error getting news source file paths"))

	req, err := http.NewRequest("GET", "/news?sources=bbc_news", nil)
	assert.Nil(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Internal Server Error")
}

func TestNewsHandlerInvalidMethod(t *testing.T) {
	handler, _, _ := setupNewsHandlerTest()

	req, err := http.NewRequest("POST", "/news", nil)
	assert.Nil(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected status Method Not Allowed")
}
