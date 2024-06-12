package entity

// ResourceName represents the name of a news resource.
type ResourceName string

// PathToFile represents the path to a file containing news information.
type PathToFile string

// Resource represents a structure containing information about a news resource,
// including its name, path to file, and source type.
type Resource struct {
	Name       ResourceName
	PathToFile PathToFile
}

func GetResourceNames(resources []Resource) []string {
	var names []string
	for _, resource := range resources {
		names = append(names, string(resource.Name))
	}
	return names
}
