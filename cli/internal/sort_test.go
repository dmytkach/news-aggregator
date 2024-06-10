package internal

import (
	"NewsAggregator/cli/internal/entity"
	"reflect"
	"testing"
	"time"
)

func TestSort(t *testing.T) {
	news := []entity.News{
		{
			Title:       "NBC_example",
			Description: "test",
			Link:        "test",
			Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
			Source:      "NBC",
		},
		{
			Title:       "CNN_example",
			Description: "test",
			Link:        "test",
			Date:        time.Date(2024, time.May, 12, 22, 30, 16, 0, time.UTC),
			Source:      "CNN",
		},
		{
			Title:       "BBC_example",
			Description: "test",
			Link:        "test",
			Date:        time.Date(2024, time.May, 15, 22, 30, 16, 0, time.UTC),
			Source:      "BBC"},
	}

	tests := []struct {
		name      string
		criterion string
		order     string
		expected  []entity.News
	}{
		{
			name:      "Sort by date ASC",
			criterion: "date",
			order:     "ASC",
			expected: []entity.News{

				{
					Title:       "CNN_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 12, 22, 30, 16, 0, time.UTC),
					Source:      "CNN",
				},
				{
					Title:       "BBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 15, 22, 30, 16, 0, time.UTC),
					Source:      "BBC",
				},
				{
					Title:       "NBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
					Source:      "NBC",
				},
			},
		},
		{
			name:      "Sort by date DESC",
			criterion: "date",
			order:     "DESC",
			expected: []entity.News{
				{
					Title:       "NBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
					Source:      "NBC",
				},
				{
					Title:       "BBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 15, 22, 30, 16, 0, time.UTC),
					Source:      "BBC",
				},
				{
					Title:       "CNN_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 12, 22, 30, 16, 0, time.UTC),
					Source:      "CNN",
				},
			},
		},
		{
			name:      "Sort by source ASC",
			criterion: "source",
			order:     "ASC",
			expected: []entity.News{
				{
					Title:       "BBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 15, 22, 30, 16, 0, time.UTC),
					Source:      "BBC",
				},
				{
					Title:       "CNN_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 12, 22, 30, 16, 0, time.UTC),
					Source:      "CNN",
				},
				{
					Title:       "NBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
					Source:      "NBC",
				},
			},
		},
		{
			name:      "Sort by source DESC",
			criterion: "source",
			order:     "DESC",
			expected: []entity.News{
				{
					Title:       "NBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
					Source:      "NBC",
				},
				{
					Title:       "CNN_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 12, 22, 30, 16, 0, time.UTC),
					Source:      "CNN",
				},
				{
					Title:       "BBC_example",
					Description: "test",
					Link:        "test",
					Date:        time.Date(2024, time.May, 15, 22, 30, 16, 0, time.UTC),
					Source:      "BBC",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Sort(news, test.criterion, test.order)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("unexpected result - got: %+v, want: %+v", result, test.expected)
			}
		})
	}
}
