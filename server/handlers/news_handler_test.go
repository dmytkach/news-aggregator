package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mockInitialize(name string) (map[string][]string, error) {
	return map[string][]string{
		"bbc_news": {
			"../../internal/testdata/bbc_news/ready_news.json",
		},
	}, nil
}

func setupResource() {
	response = httptest.NewRecorder()

	SourceInitializer = mockInitialize
}

func TestNewsHandler(t *testing.T) {
	setupResource()

	req, err := http.NewRequest("GET", "/news?sources=bbc_news&keywords=England", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(News)

	handler.ServeHTTP(rr, req)

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
	assert.Equal(t, expected, actual, "Expected response body to match")
}

func TestNewsHandlerInvalidMethod(t *testing.T) {
	setupResource()

	req, err := http.NewRequest("POST", "/news", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(News)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "Expected status Method Not Allowed")
}

func TestNewsHandlerInvalidSource(t *testing.T) {
	setupResource()

	req, err := http.NewRequest("GET", "/news?sources=invalid_source", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(News)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code, "Expected status Bad Request")
}
