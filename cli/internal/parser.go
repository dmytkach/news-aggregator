package internal

import (
	"fmt"
	"news-aggregator/internal/entity"
	"news-aggregator/internal/parser"
	"path/filepath"
	"strings"
)

// GetFileParser returns the appropriate parser implementation based on the path to file.
func GetFileParser(path entity.PathToFile) (Parser, error) {
	ext := strings.ToLower(filepath.Ext(string(path)))

	parsers := []Parser{
		&parser.Rss{FilePath: path},
		&parser.Json{FilePath: path},
		&parser.UsaToday{FilePath: path},
	}

	for _, p := range parsers {
		if p.CanParseFileType(ext) {
			return p, nil
		}
	}
	return nil, fmt.Errorf("unsupported file type: %s", ext)
}
