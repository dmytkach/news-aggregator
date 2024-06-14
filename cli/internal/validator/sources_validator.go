package validator

import (
	"strings"
)

// sourceValidator checks the sources field
type sourceValidator struct {
	baseValidator
	availableSources map[string]string
	sources          string
}

func (v sourceValidator) Validate() bool {
	sourcesList := strings.Split(v.sources, ",")
	if len(sourcesList) == 0 {
		println("Please provide at least one source using the --sources flag.")
		return false
	}
	for _, source := range sourcesList {
		if _, ok := v.availableSources[strings.ToLower(source)]; !ok {
			println("Error - Source does not exist:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}
