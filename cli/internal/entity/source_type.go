package entity

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SourceType represents the type of news source.
type SourceType interface {
	extension() string
}
type RssType struct{}

func (rt RssType) extension() string {
	return ".xml"
}

type JsonType struct{}

func (jt JsonType) extension() string {
	return ".json"
}

type HtmlType struct{}

func (ht HtmlType) extension() string {
	return ".html"
}

// AnalyzeResourceType analyzes the resource type based on the file extension.
// It returns an instance of the appropriate SourceType for the given file.
func AnalyzeResourceType(file PathToFile) SourceType {
	ext := strings.ToLower(filepath.Ext(string(file)))
	typeMap := map[string]SourceType{
		RssType{}.extension():  RssType{},
		JsonType{}.extension(): JsonType{},
		HtmlType{}.extension(): HtmlType{},
	}
	typeSource, exist := typeMap[ext]
	if !exist {
		fmt.Println("Wrong Source", ext)
		return nil
	}
	return typeSource
}
