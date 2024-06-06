package main

import (
	"NewsAggregator/cli/internal"
	"reflect"
	"testing"
)

func Test_processDateEnd(t *testing.T) {
	type args struct {
		dateEnd string
	}
	tests := []struct {
		name    string
		args    args
		want    internal.NewsFilter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processDateEnd(tt.args.dateEnd)
			if (err != nil) != tt.wantErr {
				t.Errorf("processDateEnd() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processDateEnd() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processDateStart(t *testing.T) {
	type args struct {
		dateStart string
	}
	tests := []struct {
		name    string
		args    args
		want    internal.NewsFilter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := processDateStart(tt.args.dateStart)
			if (err != nil) != tt.wantErr {
				t.Errorf("processDateStart() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processDateStart() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_processKeywords(t *testing.T) {
	type args struct {
		keywords string
	}
	tests := []struct {
		name string
		args args
		want internal.NewsFilter
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processKeywords(tt.args.keywords); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("processKeywords() = %v, want %v", got, tt.want)
			}
		})
	}
}
