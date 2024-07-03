package parser

import (
	"fmt"
	"news-aggregator/internal/entity"
	"path/filepath"
	"strings"
)

// Parser provides an API for a news parser capable of processing a specific file type.
type Parser interface {
	CanParseFileType(ext string) bool
	Parse() (entity.Feed, error)
}

// GetFileParser returns the appropriate parser implementation based on the path to file.
func GetFileParser(path entity.PathToFile) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(string(path)))

	parsers := []Parser{
		&Rss{FilePath: path},
		&Json{FilePath: path},
		&UsaToday{FilePath: path},
	}
	for p := range parsers {
		if parsers[p].CanParseFileType(ext) {
			return parsers[p], nil
		}
	}
	return nil, fmt.Errorf("unsupported file type: %s", ext)
}
