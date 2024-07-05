package service

import (
	"github.com/stretchr/testify/assert"
	"news-aggregator/internal/entity"
	"news-aggregator/server/managers"
	"os"
	"testing"
)

var url = "http://feeds.bbci.co.uk/news/rss.xml"

func TestAddSource(t *testing.T) {
	_ = os.Remove("test-resources.json")
	managers.PathToResources = "test-resources.json"
	expectedSource := entity.Source{Name: "bbc_news", PathsToFile: []entity.PathToFile{entity.PathToFile(url)}}

	source, err := AddSource(url)
	if err != nil {
		t.Errorf("Unexpected error adding source: %v", err)
	}
	if source.Name != expectedSource.Name {
		t.Errorf("Expected source name %s, got %s", expectedSource.Name, source.Name)
	}
	if len(source.PathsToFile) != len(expectedSource.PathsToFile) {
		t.Errorf("Expected %d paths to file, got %d", len(expectedSource.PathsToFile), len(source.PathsToFile))
	}

	urlWithError := "http://feeds.bbci.co.uk/news/rss.xml"
	_, err = AddSource(urlWithError)
	if err == nil {
		t.Error("Expected error but got nil")
	}

}

func TestGetSources(t *testing.T) {
	managers.PathToResources = "test-resources.json"
	_, err := AddSource(url)

	sources, err := GetSources()
	assert.NoError(t, err)
	assert.Len(t, sources, 1)
	assert.Equal(t, "bbc_news", string(sources[0].Name))
	assert.Equal(t, url, string(sources[0].PathsToFile[0]))
}

func TestGetSource(t *testing.T) {
	managers.PathToResources = "test-resources.json"
	_, err := AddSource(url)
	source, err := GetSource("bbc_news")
	assert.NoError(t, err)
	assert.Equal(t, "bbc_news", string(source.Name))
	assert.Equal(t, url, string(source.PathsToFile[0]))

	_, err = GetSource("NonExistingSource")
	assert.Error(t, err)
	assert.Equal(t, "No resources found for name: NonExistingSource", err.Error())
}

func TestUpdateSource(t *testing.T) {
	managers.PathToResources = "test-resources.json"

	_, err := AddSource(url)
	err = UpdateSource(url, "https://feeds.bbci.co.uk/news/world/asia/rss.xml")
	assert.NoError(t, err)

	err = UpdateSource("old-url", "updated-url")
	assert.Error(t, err)
	assert.Equal(t, "source with URL old-url not found", err.Error())
}

func TestRemoveSource(t *testing.T) {
	managers.PathToResources = "test-resources.json"
	_, err := AddSource(url)
	err = RemoveSource("bbc_news")
	assert.NoError(t, err)
}
