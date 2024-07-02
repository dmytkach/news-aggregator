package validator

import (
	"log"
	"time"
)

// dateEndValidator checks the end date format.
type dateEndValidator struct {
	baseValidator
	dateEnd string
}

const DateFormat = "2006-01-02"

// Validate checks that the end date is in the correct format.
func (d dateEndValidator) Validate() bool {
	if d.dateEnd != "" {
		_, err := time.Parse(DateFormat, d.dateEnd)
		if err != nil {
			log.Print("Invalid end date format. Please use YYYY-MM-DD.")
			return false
		}
	}
	return d.baseValidator.Validate()
}
