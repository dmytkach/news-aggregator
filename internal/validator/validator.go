package validator

import "news-aggregator/internal/sort"

// validator stores all the necessary parameters for performing validation.
type validator struct {
	Sources          string
	AvailableSources []string
	DateStart        string
	DateEnd          string
	Criterion        string
	Order            string
}

// Config contains all the parameters required to initialize a validator.
type Config struct {
	Sources          string
	AvailableSources []string
	DateStart        string
	DateEnd          string
	SortOptions      sort.Options
}

// NewValidator instance configured with the specified parameters.
func NewValidator(config Config) ValidatingComponent {
	return &validator{
		Sources:          config.Sources,
		AvailableSources: config.AvailableSources,
		DateStart:        config.DateStart,
		DateEnd:          config.DateEnd,
		Criterion:        config.SortOptions.Criterion,
		Order:            config.SortOptions.Order,
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
