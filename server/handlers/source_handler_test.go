package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"news-aggregator/internal/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mockAddSource(url string) (entity.Source, error) {
	return entity.Source{Name: "mockedSource", PathsToFile: []entity.PathToFile{entity.PathToFile(url)}}, nil
}

func mockFetchNews() error {
	return nil
}

func mockGetSources() ([]entity.Source, error) {
	return []entity.Source{
		{Name: "mockedSource1", PathsToFile: []entity.PathToFile{"http://example.com/rss1"}},
		{Name: "mockedSource2", PathsToFile: []entity.PathToFile{"http://example.com/rss2"}},
	}, nil
}

func mockGetSource(name string) (entity.Source, error) {
	return entity.Source{Name: entity.SourceName(name), PathsToFile: []entity.PathToFile{"http://example.com/rss"}}, nil
}

func mockUpdateSource(oldUrl, newUrl string) error {
	return nil
}

func mockRemoveSource(name string) error {
	return nil
}

var (
	response      *httptest.ResponseRecorder
	mockSourceURL = "http://feeds.bbci.co.uk/news/rss.xml"
)

func setup() {
	response = httptest.NewRecorder()

	addSourceFunc = mockAddSource
	fetchNewsFunc = mockFetchNews
	getSourcesFunc = mockGetSources
	getSourceFunc = mockGetSource
	updateSourceFunc = mockUpdateSource
	removeSourceFunc = mockRemoveSource
}
func TestDownloadSource(t *testing.T) {
	setup()

	req, err := http.NewRequest("POST", "/sources?url="+mockSourceURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(Sources)
	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var result entity.Source
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, "mockedSource", string(result.Name))
	assert.Equal(t, []entity.PathToFile{entity.PathToFile(mockSourceURL)}, result.PathsToFile)
}

func TestGetSources(t *testing.T) {
	setup()

	req, err := http.NewRequest("GET", "/sources", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(Sources)
	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)

	var result []entity.Source
	err = json.NewDecoder(response.Body).Decode(&result)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, "mockedSource1", string(result[0].Name))
	assert.Equal(t, "mockedSource2", string(result[1].Name))
}

func TestUpdateSource(t *testing.T) {
	setup()

	req, err := http.NewRequest("PUT", "/sources?oldUrl=old&newUrl=new", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(Sources)
	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
}

func TestRemoveSource(t *testing.T) {
	setup()

	req, err := http.NewRequest("DELETE", "/sources?name=example", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler := http.HandlerFunc(Sources)
	handler.ServeHTTP(response, req)

	assert.Equal(t, http.StatusOK, response.Code)
}
