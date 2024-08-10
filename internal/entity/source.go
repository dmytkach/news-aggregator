package entity

// SourceName represents the name of a news source.
type SourceName string

type PathToFile string

type Source struct {
	Name       SourceName
	PathToFile PathToFile
}
