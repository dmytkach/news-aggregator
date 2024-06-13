package internal

import (
	"fmt"
	"news-aggregator/cli/internal/entity"
	"news-aggregator/cli/internal/parser"
)

// Parser is an interface that defines the Parse method for assembling news from a given file.
type Parser interface {
	Parse() ([]entity.News, error)
}

// New returns the appropriate parser implementation based on the type of source.
func New(Path entity.PathToFile) Parser {
	typeOfSource := entity.AnalyzeResourceType(Path)
	parserMap := map[entity.SourceType]Parser{
		entity.RssType{}:  &parser.Rss{FilePath: Path},
		entity.JsonType{}: &parser.Json{FilePath: Path},
		entity.HtmlType{}: &parser.UsaToday{FilePath: Path},
	}
	p, exist := parserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}

	return p
}
