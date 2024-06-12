package internal

import (
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/filters"
	"reflect"
	"testing"
	"time"
)

func TestAggregate(t *testing.T) {
	type args struct {
		Resources   []entity.Resource
		Sources     []string
		NewsFilters []NewsFilter
	}
	tests := []struct {
		name string
		args args
		want []entity.News
	}{{
		name: "should aggregate news on given Sources applying NewsFilters.",
		args: args{
			Resources: []entity.Resource{
				{
					"RSS", "testdata/news.xml",
				},
				{
					"HTML", "testdata/news.html",
				},
				{
					"JSON", "testdata/news.json",
				},
			},
			Sources: []string{"RSS"},
			NewsFilters: []NewsFilter{
				&filters.Keyword{Keywords: []string{"Dymond"}},
				&filters.DateStart{StartDate: time.Date(2024, time.May, 17, 10, 10, 10, 0, time.UTC)},
				&filters.DateEnd{EndDate: time.Date(2024, time.May, 18, 23, 30, 10, 0, time.UTC)},
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
			a := Aggregator{
				Resources:   tt.args.Resources,
				Sources:     tt.args.Sources,
				NewsFilters: tt.args.NewsFilters,
			}
			if got := a.Aggregate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestAggregator_applyFilters(t *testing.T) {
	type fields struct {
		NewsFilters []NewsFilter
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
			NewsFilters: []NewsFilter{
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
			a := Aggregator{
				NewsFilters: tt.fields.NewsFilters,
			}
			if got := a.applyFilters(tt.args.news); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_collectNews(t *testing.T) {
	type fields struct {
		Resources []entity.Resource
	}
	type args struct {
		sources []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []entity.News
	}{
		{
			name: "should aggregate news on given Sources applying NewsFilters.",
			fields: fields{
				Resources: []entity.Resource{
					{
						"RSS", "testdata/news.xml",
					},
					{
						"HTML", "testdata/news.html",
					},
					{
						"JSON", "testdata/news.json",
					},
				},
			},
			args: args{
				sources: []string{"RSS"},
			},
			want: []entity.News{

				{
					Title:       "Boy's body found in river as second teen 'critical'",
					Description: "The 14-year-old's body was found in the River Tyne and the younger boy was rescued and taken to hospital.",
					Link:        "https://www.bbc.com/news/articles/cnee7lp7mgdo",
					Date:        time.Date(2024, time.May, 19, 12, 20, 49, 0, time.UTC),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aggregator{
				Resources: tt.fields.Resources,
			}
			if got := a.collectNews(tt.args.sources); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("collectNews() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_getNewsForSource(t *testing.T) {
	type fields struct {
		Resources []entity.Resource
	}
	type args struct {
		sourceName string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []entity.News
	}{
		{
			name: "should aggregate news on given Sources applying NewsFilters.",
			fields: fields{
				Resources: []entity.Resource{
					{
						"RSS", "testdata/news.xml",
					},
					{
						"HTML", "testdata/news.html",
					},
					{
						"JSON", "testdata/news.json",
					},
				},
			},
			args: args{
				"RSS",
			},
			want: []entity.News{
				{
					Title:       "Boy's body found in river as second teen 'critical'",
					Description: "The 14-year-old's body was found in the River Tyne and the younger boy was rescued and taken to hospital.",
					Link:        "https://www.bbc.com/news/articles/cnee7lp7mgdo",
					Date:        time.Date(2024, time.May, 19, 12, 20, 49, 0, time.UTC),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Aggregator{
				Resources: tt.fields.Resources,
			}
			if got := a.getNewsForSource(tt.args.sourceName); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getNewsForSource() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAggregator_getResourceNews(t *testing.T) {

	type args struct {
		path entity.PathToFile
	}
	tests := []struct {
		name string
		args args
		want []entity.News
	}{
		{
			name: "should aggregate news on given Sources applying NewsFilters.",
			args: args{
				"testdata/news.xml",
			},
			want: []entity.News{
				{
					Title:       "Boy's body found in river as second teen 'critical'",
					Description: "The 14-year-old's body was found in the River Tyne and the younger boy was rescued and taken to hospital.",
					Link:        "https://www.bbc.com/news/articles/cnee7lp7mgdo",
					Date:        time.Date(2024, time.May, 19, 12, 20, 49, 0, time.UTC),
				},
				{
					Title:       "Su and Steve fought for justice, but didn't live to see it",
					Description: "Su Gorman and Steve Dymond helped to expose an NHS scandal - but they died without seeing justice.",
					Link:        "https://www.bbc.co.uk/news/health-69018125",
					Date:        time.Date(2024, time.May, 18, 23, 05, 19, 0, time.UTC),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Aggregator{}
			if got := a.getResourceNews(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResourceNews() = %v, want %v", got, tt.want)
			}
		})
	}
}
