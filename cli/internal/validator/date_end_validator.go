package validator

import "time"

// DateEndValidator checks the dateEnd field
type DateEndValidator struct {
	BaseValidator
	DateEnd string
}

func (d DateEndValidator) Validate() bool {
	if d.DateEnd != "" {
		_, err := time.Parse("2006-01-02", d.DateEnd)
		if err != nil {
			println("Invalid end date format. Please use YYYY-MM-DD.")
			return false
		}
	}
	return d.BaseValidator.Validate()
}
