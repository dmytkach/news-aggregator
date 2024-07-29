package validator

import (
	"errors"
	"fmt"
	"strings"
)

// sourceValidator checks the sources field
type sourceValidator struct {
	baseValidator
	availableSources []string
	sources          string
}

// Validate checks that the given sources exist.
func (v sourceValidator) Validate() error {
	sourcesList := strings.Split(v.sources, ",")
	if len(v.availableSources) == 0 {
		//log.Println("Not found available sources.")
		return errors.New("not found available sources")
	}
	for _, source := range sourcesList {
		if strings.TrimSpace(source) == "" {
			//log.Println("Please provide at least one source.")
			return errors.New("please provide at least one source")
		}
		if !v.exist(source) {
			//log.Println("Source not available:", source)
			return fmt.Errorf("source not available: %s", source)
		}
	}
	return v.baseValidator.Validate()
}

// exist check elements existing in slice.
func (v sourceValidator) exist(source string) bool {
	for _, s := range v.availableSources {
		if strings.EqualFold(s, source) {
			return true
		}
	}
	return false
}
