package entity

// ResourceName represents the name of a news resource.
type ResourceName string

// PathToFile represents the path to a file containing news information.
type PathToFile string
type Resource struct {
	Name       ResourceName
	PathToFile PathToFile
}
