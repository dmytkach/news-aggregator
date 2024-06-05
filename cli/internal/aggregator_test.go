package internal

import (
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/filters"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func setWorkingDirectory(t *testing.T) func() {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	rootPath := filepath.Join(wd, "../..")
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

func TestAggregate(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	type args struct {
		sources []string
		filters []NewsFilter
	}
	tests := []struct {
		name string
		args args
		want []entity.News
	}{{
		name: "should work with bbc resource",
		args: args{
			sources: []string{"BBC"},
			filters: []NewsFilter{
				&filters.Keyword{Keywords: []string{"Taiwan", "Israel"}},
				&filters.DateStart{StartDate: time.Date(2024, time.May, 17, 10, 10, 10, 0, time.UTC)},
				&filters.DateEnd{EndDate: time.Date(2024, time.May, 18, 23, 30, 10, 0, time.UTC)},
			},
		},
		want: []entity.News{
			{
				Title:       "Israel war cabinet minister vows to quit if there is no post-war plan for Gaza",
				Description: "Recent weeks have seen an increasingly public rift over how Gaza should be governed after the war.",
				Link:        "https://www.bbc.com/news/articles/cekkz82gnzgo",
				Date:        time.Date(2024, time.May, 18, 23, 22, 26, 0, time.UTC),
			},
			{
				Title:       "Taiwan's steely leader rewrote the book on how to deal with China",
				Description: "Taiwan's outgoing president says boosting her country's military was the only way to defy China's threat.",
				Link:        "https://www.bbc.com/news/articles/ceklk794102o",
				Date:        time.Date(2024, time.May, 18, 22, 24, 57, 0, time.UTC),
			},
			{
				Title:       "Stories of the hostages taken by Hamas from Israel",
				Description: "It is thought more than 100 Israelis are still being held hostage in Gaza after the 7 October attacks.",
				Link:        "https://www.bbc.co.uk/news/world-middle-east-67053011",
				Date:        time.Date(2024, time.May, 18, 22, 30, 16, 0, time.UTC),
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Aggregate(tt.args.sources, tt.args.filters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Aggregate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_applyFilters(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	type args struct {
		news    []entity.News
		filters []NewsFilter
	}
	tests := []struct {
		name string
		args args
		want []entity.News
	}{{
		name: "should work with nbc resource",
		args: args{
			news: collectNews([]string{"NBC"}),
			filters: []NewsFilter{
				&filters.Keyword{Keywords: []string{"Taiwan"}},
			},
		},
		want: []entity.News{
			{
				Title:       "Taiwan lawmakers exchange blows in bitter dispute over parliament reforms",
				Description: "Lawmakers in Taiwan shoved, tackled and hit each other in a bitter dispute about parliamentary reforms, days before President-elect Lai Ching-te takes office/",
				Link:        "https://www.nbcnews.com/news/world/taiwan-lawmakers-exchange-blows-bitter-dispute-parliament-reforms-rcna152763",
				Date:        time.Date(2024, time.May, 17, 15, 42, 56, 0, time.UTC),
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := applyFilters(tt.args.news, tt.args.filters); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("applyFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_collectNews(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	type args struct {
		sources []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "should work with nbc resource",
		args: args{
			sources: []string{"usatoday"},
		},
		want: 33,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := collectNews(tt.args.sources); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("collectNews() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_getForSource(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	type args struct {
		sourceName string
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "should work with nbc resource",
		args: args{
			sourceName: "abc",
		},
		want: 25,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getForSource(tt.args.sourceName); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getForSource() = %v, want %v", len(got), tt.want)
			}
		})
	}
}

func Test_getResourceNews(t *testing.T) {
	restoreWD := setWorkingDirectory(t)
	defer restoreWD()
	type args struct {
		resource entity.Resource
	}
	tests := []struct {
		name string
		args args
		want int
	}{{
		name: "should work with nbc resource",
		args: args{
			resource: entity.Resource{Name: "washington", PathToFile: "resources/washingtontimes-world-category-19-05-24.xml"},
		},
		want: 20,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getResourceNews(tt.args.resource); !reflect.DeepEqual(len(got), tt.want) {
				t.Errorf("getResourceNews() = %v, want %v", len(got), tt.want)
			}
		})
	}
}
