package validator

import (
	"errors"
	"strings"
)

type sortOptionsValidator struct {
	baseValidator
	criterion string
	order     string
}

// Validate checks if the provided sorting criterion and order are valid.
// If a criterion is specified, it must be either "date" or "source".
// If an order is specified, it must be either "asc" or "desc".
func (v sortOptionsValidator) Validate() error {
	if v.criterion != "" {
		if !strings.EqualFold(v.criterion, "date") && !strings.EqualFold(v.criterion, "source") {
			//log.Println("Invalid sort criterion. Please use `date` or `source`")
			return errors.New("invalid sort criterion. Please use `date` or `source`")
		}
	}
	if v.order != "" {
		if !strings.EqualFold(v.order, "desc") && !strings.EqualFold(v.order, "asc") {
			//log.Println("Invalid order criterion. Please use `asc` or `desc`")
			return errors.New("invalid order criterion. Please use `asc` or `desc`")
		}
	}
	return v.baseValidator.Validate()
}
