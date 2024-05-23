package parser

import (
	"NewsAggregator/cli/internal/entity"
	"fmt"
)

// Parser is an interface that defines the Parser method for assembling news from a given file.
type Parser interface {
	Parse(FileName string) ([]entity.News, error)
}

// GetParser returns the appropriate parser implementation based on the type of source.
func GetParser(typeOfSource entity.SourceType) Parser {
	parserMap := map[entity.SourceType]Parser{
		"RSS":  &RssParser{},
		"JSON": &JsonParser{},
		"Html": &HtmlParser{},
	}
	parser, exist := parserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}

	return parser
}
