package handlers

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers/mock_managers"
	"testing"
	"time"
)

func setupNewsHandlerTest(ctrl *gomock.Controller) (*NewsHandler, *mock_managers.MockNewsManager, *mock_managers.MockSourceManager) {
	mockNewsManager := mock_managers.NewMockNewsManager(ctrl)
	mockSourceManager := mock_managers.NewMockSourceManager(ctrl)
	handler := &NewsHandler{
		NewsManager:   mockNewsManager,
		SourceManager: mockSourceManager,
	}
	return handler, mockNewsManager, mockSourceManager
}

func TestNewsHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest(ctrl)

	mockSources := []entity.Source{
		{Name: "bbc_news"},
	}
	mockSourceManager.EXPECT().GetSources().Return(mockSources, nil)

	mockNewsManager.EXPECT().GetNewsSourceFilePath([]string{"bbc_news"}).
		Return(map[string][]string{
			"bbc_news": {"../../internal/testdata/bbc_news/ready_news.json"},
		}, nil)

	req, err := http.NewRequest("GET", "/news?sources=bbc_news&keywords=England", nil)
	assert.NoError(t, err, "Expected no error creating request")

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
	assert.NoError(t, err, "Expected no error decoding response body")
	assert.ElementsMatch(t, expected, actual, "Expected response body to match")
}

func TestNewsHandlerInvalidSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest(ctrl)

	var mockSources []entity.Source
	mockSourceManager.EXPECT().GetSources().Return(mockSources, nil)

	mockNewsManager.EXPECT().GetNewsSourceFilePath(gomock.Any()).
		Return(map[string][]string{}, errors.New("error getting news source file paths"))

	req, err := http.NewRequest("GET", "/news?sources=invalid_source", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}

func TestNewsHandlerErrorGettingSources(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest(ctrl)

	// Настроим ожидание для вызова GetSources
	mockSourceManager.EXPECT().GetSources().Return(nil, errors.New("error getting sources"))

	// Ожидаем, что GetNewsSourceFilePath вызовется с любыми аргументами, чтобы не было ошибок в тесте
	mockNewsManager.EXPECT().GetNewsSourceFilePath(gomock.Any()).Return(nil, nil).Times(0)

	req, err := http.NewRequest("GET", "/news?sources=invalid_source", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code, "Expected status Internal Server Error")
}

func TestNewsHandlerErrorGettingNewsSourceFilePath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, mockNewsManager, mockSourceManager := setupNewsHandlerTest(ctrl)

	mockSources := []entity.Source{
		{Name: "bbc_news"},
	}
	mockSourceManager.EXPECT().GetSources().Return(mockSources, nil)
	mockNewsManager.EXPECT().GetNewsSourceFilePath([]string{"bbc_news"}).
		Return(map[string][]string{}, errors.New("error getting news source file paths"))

	req, err := http.NewRequest("GET", "/news?sources=bbc_news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}

func TestNewsHandlerInvalidMethod(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	handler, _, _ := setupNewsHandlerTest(ctrl)

	req, err := http.NewRequest("POST", "/news", nil)
	assert.NoError(t, err, "Expected no error creating request")

	rr := httptest.NewRecorder()
	httpHandler := http.HandlerFunc(handler.News)

	httpHandler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected status Method Not Allowed")
}
