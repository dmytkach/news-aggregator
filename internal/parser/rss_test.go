package parser

import (
	"news-aggregator/internal/entity"
	"reflect"
	"testing"
	"time"
)

// Unit test for rss parser.
func TestRss_Parse(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []entity.News
		wantErr bool
	}{
		{
			name: "Test parsing valid RSS file",
			file: "../testdata/news.xml",
			want: []entity.News{
				{
					Title:       "Boy's body found in river as second teen 'critical'",
					Description: "The 14-year-old's body was found in the River Tyne and the younger boy was rescued and taken to hospital.",
					Link:        "https://www.bbc.com/news/articles/cnee7lp7mgdo",
					Date:        time.Date(2024, 5, 19, 12, 20, 49, 0, time.UTC),
					Source:      "BBC News",
				},
				{
					Title:       "Su and Steve fought for justice, but didn't live to see it",
					Description: "Su Gorman and Steve Dymond helped to expose an NHS scandal - but they died without seeing justice.",
					Link:        "https://www.bbc.co.uk/news/health-69018125",
					Date:        time.Date(2024, 5, 18, 23, 5, 19, 0, time.UTC),
					Source:      "BBC News",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test parsing invalid RSS file",
			file:    "../testdata/invalid_news.xml",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Test parsing missing RSS file",
			file:    "../testdata/nonexistent_file.xml",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rssParser := &Rss{
				FilePath: entity.PathToFile(tt.file),
			}
			got, err := rssParser.Parse()
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
