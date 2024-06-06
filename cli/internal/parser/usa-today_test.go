package parser

import (
	"NewsAggregator/cli/internal/entity"
	"reflect"
	"testing"
)

func TestUsaToday_Parse(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []entity.News
		wantErr bool
	}{{
		name:    "should collect news from the specified sources.",
		file:    "../testdata/news.xml",
		want:    nil,
		wantErr: false,
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usaTodayParser := &UsaToday{
				FilePath: entity.PathToFile(tt.file),
			}
			got, err := usaTodayParser.Parse()
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
