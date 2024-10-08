package validator

import (
	"errors"
	"time"
)

// dateEndValidator checks the end date format.
type dateEndValidator struct {
	baseValidator
	dateEnd string
}

const DateFormat = "2006-01-02"

// Validate checks that the end date is in the correct format
// and the news is not earlier than 1900
func (d dateEndValidator) Validate() error {
	if d.dateEnd != "" {
		date, err := time.Parse(DateFormat, d.dateEnd)
		if err != nil {
			//log.Println("Invalid end date format. Please use YYYY-MM-DD.")
			return errors.New("invalid end date format. Please use YYYY-MM-DD")
		}
		if date.Before(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)) {
			//log.Println("News for this period is not available.")
			return errors.New("news for this period is not available")
		}
	}
	return d.baseValidator.Validate()
}
