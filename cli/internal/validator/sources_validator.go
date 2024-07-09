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
	for _, source := range sourcesList {
		if source == "" {
			println("Please provide source using the --sources flag.")
			return false
		}
		if _, ok := v.availableSources[strings.ToLower(source)]; !ok {
			println("Error - Source does not exist:", source)
			return false
		}
	}
	return v.baseValidator.Validate()
}
