package entity

import (
	"fmt"
	"path/filepath"
	"strings"
)

// SourceType represents the type of news source.
type SourceType interface {
	TypeName() string
}
type RssType struct{}

func (rt RssType) TypeName() string {
	return ".xml"
}

type JsonType struct{}

func (jt JsonType) TypeName() string {
	return ".json"
}

type HtmlType struct{}

func (ht HtmlType) TypeName() string {
	return ".html"
}
func AnalyzeResourceType(file PathToFile) SourceType {
	ext := strings.ToLower(filepath.Ext(string(file)))
	sourceMap := map[string]SourceType{
		RssType{}.TypeName():  RssType{},
		JsonType{}.TypeName(): JsonType{},
		HtmlType{}.TypeName(): HtmlType{},
	}
	typeSource, exist := sourceMap[ext]
	if !exist {
		fmt.Println("Wrong Source", ext)
		return nil
	}

	return typeSource
}
