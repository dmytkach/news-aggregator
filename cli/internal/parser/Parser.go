package parser

import (
	"NewsAggregator/cli/internal/entity"
	"fmt"
)

// Parser is an interface that defines the Parser method for assembling news from a given file.
type Parser interface {
	Parse() ([]entity.News, error)
}

// New returns the appropriate parser implementation based on the type of source.
func New(path entity.PathToFile) Parser {
	typeOfSource := entity.AnalyzeResourceType(path)
	parserMap := map[entity.SourceType]Parser{
		entity.RssType{}:  &RssParser{path},
		entity.JsonType{}: &JsonParser{path},
		entity.HtmlType{}: &HtmlParser{path},
	}
	parser, exist := parserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}

	return parser
}
