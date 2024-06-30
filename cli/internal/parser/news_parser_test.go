package parser

import (
	"news-aggregator/internal/entity"
	"reflect"
	"testing"
)

func TestNewsParser_CanParseFileType(t *testing.T) {
	type fields struct {
		FilePath entity.PathToFile
	}
	type args struct {
		ext string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsParser := &NewsParser{
				FilePath: tt.fields.FilePath,
			}
			if got := newsParser.CanParseFileType(tt.args.ext); got != tt.want {
				t.Errorf("CanParseFileType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewsParser_Parse(t *testing.T) {
	type fields struct {
		FilePath entity.PathToFile
	}
	tests := []struct {
		name    string
		fields  fields
		want    []entity.News
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newsParser := &NewsParser{
				FilePath: tt.fields.FilePath,
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
