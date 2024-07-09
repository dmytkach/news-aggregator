package validator

import "news-aggregator/internal/sort"

type validator struct {
	Sources          string
	AvailableSources map[string]string
	DateStart        string
	DateEnd          string
	Criterion        string
	Order            string
}

func NewValidator(sources string, availableSources map[string]string, dateStart, dateEnd string, sort sort.Options) ValidatingComponent {
	return &validator{
		Sources:          sources,
		AvailableSources: availableSources,
		DateStart:        dateStart,
		DateEnd:          dateEnd,
		Criterion:        sort.Criterion,
		Order:            sort.Order,
	}
}

// ValidatingComponent defines the contract for components responsible for performing validation.
type ValidatingComponent interface {
	Validate() bool
}

// Validate all implementations with the chain of responsibility pattern.
func (v validator) Validate() bool {
	sourceValidator := &sourceValidator{sources: v.Sources, availableSources: v.AvailableSources}
	dateStartValidator := &dateStartValidator{dateStart: v.DateStart}
	dateEndValidator := &dateEndValidator{dateEnd: v.DateEnd}
	sortOptionsValidator := &sortOptionsValidator{criterion: v.Criterion, order: v.Order}

	sourceValidator.SetNext(dateStartValidator)
	dateStartValidator.SetNext(dateEndValidator)
	dateEndValidator.SetNext(sortOptionsValidator)

	return sourceValidator.Validate()
}
