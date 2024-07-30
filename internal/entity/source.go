package entity

// SourceName represents the name of a news source.
type SourceName string

// PathToFile represents the path to a file containing news information.
type PathToFile string
type Source struct {
	Name        SourceName
	PathsToFile []PathToFile
}
