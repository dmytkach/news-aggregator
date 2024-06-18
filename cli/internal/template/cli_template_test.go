package template_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/wk8/go-ordered-map"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/template"
)

// Mock data and functions
var testNews = []entity.News{
	{Title: "President speaks", Description: "The president gave a speech", Source: "BBC"},
	{Title: "New law signed", Description: "The president signed a new law", Source: "NBC"},
	{Title: "President travels", Description: "The president is traveling", Source: "BBC"},
}

func setWorkingDirectory(t *testing.T) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	rootPath := filepath.Join(wd, "../../..")
	err = os.Chdir(rootPath)
	if err != nil {
		t.Fatalf("failed to change working directory: %v", err)
	}
	return func() {
		err := os.Chdir(wd)
		if err != nil {
			t.Fatalf("failed to restore working directory: %v", err)
		}
	}
}
func TestCreate(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	tests := []struct {
		name     string
		keywords string
		wantErr  bool
	}{
		{"NoKeywords", "", false},
		{"WithKeywords", "breaking,world", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := template.Data{}
			tmpl, err := data.Create(tt.keywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tmpl == nil && !tt.wantErr {
				t.Error("Create() got = nil, want non-nil template")
			}
		})
	}
}
func TestGroupNews(t *testing.T) {
	news := testNews

	grouped := group(news)
	expectedGroups := 2

	if grouped.Len() != expectedGroups {
		t.Errorf("expected %d groups, got %d", expectedGroups, grouped.Len())
	}

	bbcNews, _ := grouped.Get("BBC")
	if len(bbcNews.([]entity.News)) != 2 {
		t.Errorf("expected 2 news items for BBC, got %d", len(bbcNews.([]entity.News)))
	}
}

func TestPrepare(t *testing.T) {
	data := template.Data{
		News: testNews,
	}

	preparedData := data.Prepare()
	expectedGroups := 2

	if len(preparedData.Grouped) != expectedGroups {
		t.Errorf("expected %d groups, got %d", expectedGroups, len(preparedData.Grouped))
	}

	bbcNews := preparedData.Grouped[0]
	if bbcNews.Source != "BBC" {
		t.Errorf("expected source to be BBC, got %s", bbcNews.Source)
	}

	if len(bbcNews.NewsList) != 2 {
		t.Errorf("expected 2 news items for BBC, got %d", len(bbcNews.NewsList))
	}
}

// Mock function for grouping news items by their source.
func group(news []entity.News) *orderedmap.OrderedMap {
	grouped := orderedmap.New()
	for _, item := range news {
		if value, ok := grouped.Get(item.Source); ok {
			grouped.Set(item.Source, append(value.([]entity.News), item))
		} else {
			grouped.Set(item.Source, []entity.News{item})
		}
	}
	return grouped
}
