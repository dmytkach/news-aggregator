package parser

import (
	"news-aggregator/internal/entity"
	"reflect"
	"testing"
	"time"
)

func TestNewsParser_CanParseFileType(t *testing.T) {
	tests := []struct {
		name string
		ext  string
		want bool
	}{
		{
			name: "JSON file",
			ext:  ".json",
			want: true,
		},
		{
			name: "Other file types",
			ext:  ".xml",
			want: false,
		},
		{
			name: "No extension",
			ext:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsParser := &NewsParser{}
			if got := newsParser.CanParseFileType(tt.ext); got != tt.want {
				t.Errorf("CanParseFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestNewsParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []entity.News
		wantErr bool
	}{
		{
			name: "Test parsing valid JSON file with ready pre-processed news",
			file: "../testdata/ready_news.json",
			want: []entity.News{
				{
					Title:       "Why parents are locking themselves in cells at Korean 'happiness factory'",
					Description: "Some South Koreans are spending time in a cell to try to understand their socially isolated children.",
					Link:        "https://www.bbc.com/news/articles/c2x0le06kn7o",
					Date:        time.Date(2024, 6, 30, 1, 6, 19, 0, time.UTC),
					Source:      "bbc news",
				},
				{
					Title:       "Watch England fans go wild as Bellingham scores late equaliser",
					Description: "Watch England fans erupt at a fanpark in Wembley as Jude Bellingham scores a late stunner to send the game to extra time against Slovakia in the Euro 2024 last-16 match in Gelsenkirchen.\n\n\n",
					Link:        "https://www.bbc.com/sport/football/videos/cl4yj1ve5z7o",
					Date:        time.Date(2024, 6, 30, 19, 31, 26, 0, time.UTC),
					Source:      "bbc news",
				},
				{
					Title:       "Sunak and Labour's Final Sunday Pitch",
					Description: "The PM and Labour’s campaign chief are in Laura’s BBC One studio",
					Link:        "https://www.bbc.co.uk/sounds/play/p0j7g5dz",
					Date:        time.Date(2024, 6, 30, 13, 52, 0, 0, time.UTC),
					Source:      "bbc news",
				},
			},
			wantErr: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsParser := &NewsParser{
				FilePath: entity.PathToFile(tt.file),
			}
			got, err := newsParser.Parse()
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
