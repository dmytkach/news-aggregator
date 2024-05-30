package agregator

import (
	"NewsAggregator/cli/internal/entity"
	parser2 "NewsAggregator/cli/internal/parser"
	"fmt"
)

// Parser is an interface that defines the Parser method for assembling news from a given file.
type Parser interface {
	Parse() ([]entity.News, error)
}

// New returns the appropriate parser implementation based on the type of source.
func New(Path entity.PathToFile) Parser {
	typeOfSource := entity.AnalyzeResourceType(Path)
	parserMap := map[entity.SourceType]Parser{
		entity.RssType{}:  &parser2.RssParser{FilePath: Path},
		entity.JsonType{}: &parser2.JsonParser{FilePath: Path},
		entity.HtmlType{}: &parser2.UsaTodayParser{FilePath: Path},
	}
	parser, exist := parserMap[typeOfSource]
	if !exist {
		fmt.Println("Wrong Source", typeOfSource)
		return nil
	}

	return parser
}
