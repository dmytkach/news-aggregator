package aggregator

import (
	"NewsAggregator/cli/internal/entity"
	"NewsAggregator/cli/internal/parser"
	"fmt"
)

// Parser is an interface that defines the Parse method for assembling news from a given file.
type Parser interface {
	Parse() ([]entity.News, error)
}

// New returns the appropriate parser implementation based on the type of source.
func New(Path entity.PathToFile) Parser {
	typeOfSource := entity.AnalyzeResourceType(Path)
	parserMap := map[entity.SourceType]Parser{
		entity.RssType{}:  &parser.RssParser{FilePath: Path},
		entity.JsonType{}: &parser.JsonParser{FilePath: Path},
		entity.HtmlType{}: &parser.UsaTodayParser{FilePath: Path},
	}
	p, exist := parserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}

	return p
}
