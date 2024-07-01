package parser

import (
	"news-aggregator/internal/entity"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		Path entity.PathToFile
	}
	tests := []struct {
		name string
		args args
		want []Parser
	}{
		{
			name: "RSS file",
			args: args{Path: "testdata/news.xml"},
			want: []Parser{&Rss{FilePath: "testdata/news.xml"}},
		},
		{
			name: "JSON file",
			args: args{Path: "testdata/news.json"},
			want: []Parser{
				&Json{FilePath: "testdata/news.json"},
				&NewsParser{FilePath: "testdata/news.json"}},
		},
		{
			name: "HTML file",
			args: args{Path: "testdata/news.html"},
			want: []Parser{&UsaToday{FilePath: "testdata/news.html"}},
		},
		{
			name: "Unsupported file type",
			args: args{Path: "testdata/news.txt"},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := GetFileParser(tt.args.Path)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetFileParser() = %v, want %v", got, tt.want)
			}
		})
	}
}
