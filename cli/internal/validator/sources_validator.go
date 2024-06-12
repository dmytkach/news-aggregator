package validator

import (
	"strings"
)

// SourceValidator checks the sources field
type SourceValidator struct {
	BaseValidator
	AvailableSources []string
	Sources          string
}

func (v SourceValidator) Validate() bool {
	sourcesList := strings.Split(v.Sources, ",")
	if len(sourcesList) == 0 {
		println("Please provide at least one source using the --sources flag.")
		return false
	}
	for _, source := range sourcesList {
		found := false
		for i := range v.AvailableSources {
			if strings.EqualFold(source, v.AvailableSources[i]) {
				found = true
				break
			}
		}
		if !found {
			println("Error - Source does not exist:", source)
			return false
		}
	}
	return v.BaseValidator.Validate()
}
