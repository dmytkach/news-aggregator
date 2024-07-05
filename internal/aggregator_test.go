package internal

import (
	"news-aggregator/internal/entity"
	"news-aggregator/internal/filters"
	"news-aggregator/internal/initializers"
	"reflect"
	"testing"
	"time"
)

func TestAggregate(t *testing.T) {
	type args struct {
		News        map[string][]string
		Sources     string
		NewsFilters []filters.NewsFilter
	}
	resources, err := initializers.LoadSources("testdata/")
	if err != nil {
		t.Errorf("error with load resource from Folder")
	}
	tests := []struct {
		name string
		args args
		want []entity.News
	}{{
		name: "should aggregate news on given Sources applying NewsFilters.",
		args: args{
			News:    resources,
			Sources: "bbc",
			NewsFilters: []filters.NewsFilter{
				&filters.Keyword{Keywords: []string{"South"}},
				&filters.DateStart{StartDate: time.Date(2024, time.May, 17, 10, 10, 10, 0, time.UTC)},
				&filters.DateEnd{EndDate: time.Date(2024, time.July, 18, 23, 30, 10, 0, time.UTC)},
			},
		},
		want: []entity.News{
			{
				Title:       "Why parents are locking themselves in cells at Korean 'happiness factory'",
				Description: "Some South Koreans are spending time in a cell to try to understand their socially isolated children.",
				Link:        "https://www.bbc.com/news/articles/c2x0le06kn7o",
				Date:        time.Date(2024, time.June, 30, 1, 6, 19, 0, time.UTC),
				Source:      "bbc_news",
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				Resources:   tt.args.News,
				Sources:     tt.args.Sources,
				NewsFilters: tt.args.NewsFilters,
			}
			if got, _ := a.Aggregate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAggregator_applyFilters(t *testing.T) {
	type fields struct {
		Sources     string
		NewsFilters []filters.NewsFilter
	}
	type args struct {
		news []entity.News
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []entity.News
	}{{
		name: "should aggregate news on given Sources applying NewsFilters.",
		fields: fields{
			NewsFilters: []filters.NewsFilter{
				&filters.Keyword{Keywords: []string{"Dymond"}},
				&filters.DateStart{StartDate: time.Date(2024, time.May, 17, 10, 10, 10, 0, time.UTC)},
				&filters.DateEnd{EndDate: time.Date(2024, time.May, 18, 23, 30, 10, 0, time.UTC)},
			},
		},
		args: args{
			news: []entity.News{
				{
					Title:       "Boy's body found in river as second teen 'critical'",
					Description: "The 14-year-old's body was found in the River Tyne and the younger boy was rescued and taken to hospital.",
					Link:        "https://www.bbc.com/news/articles/cnee7lp7mgdo",
					Date:        time.Date(2024, time.May, 18, 23, 05, 19, 0, time.UTC),
					Source:      "RSS",
				},
				{
					Title:       "Su and Steve fought for justice, but didn't live to see it",
					Description: "Su Gorman and Steve Dymond helped to expose an NHS scandal - but they died without seeing justice.",
					Link:        "https://www.bbc.co.uk/news/health-69018125",
					Date:        time.Date(2024, time.May, 18, 23, 05, 19, 0, time.UTC),
					Source:      "RSS",
				},
			},
		},
		want: []entity.News{
			{
				Title:       "Su and Steve fought for justice, but didn't live to see it",
				Description: "Su Gorman and Steve Dymond helped to expose an NHS scandal - but they died without seeing justice.",
				Link:        "https://www.bbc.co.uk/news/health-69018125",
				Date:        time.Date(2024, time.May, 18, 23, 05, 19, 0, time.UTC),
				Source:      "RSS",
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &aggregator{
				Sources:     tt.fields.Sources,
				NewsFilters: tt.fields.NewsFilters,
			}
			if got := a.applyFilters(tt.args.news); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
