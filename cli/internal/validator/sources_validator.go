package validator

import (
	"strings"
)

// sourceValidator checks the sources field
type sourceValidator struct {
	baseValidator
	availableSources []string
	sources          string
}

func (v sourceValidator) Validate() bool {
	sourcesList := strings.Split(v.sources, ",")
	if len(sourcesList) == 0 {
		println("Please provide at least one source using the --sources flag.")
		return false
	}
	for _, source := range sourcesList {
		found := false
		for i := range v.availableSources {
			if strings.EqualFold(source, v.availableSources[i]) {
				found = true
				break
			}
		}
		if !found {
			println("Error - Source does not exist:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}
