package validator

// baseValidator provides a basic implementation of ValidatingComponent,
// supporting Chain of Responsibility.
type baseValidator struct {
	next ValidatingComponent
}

// SetNext Validator in Chain.
func (b *baseValidator) SetNext(next ValidatingComponent) {
	b.next = next
}

// Validate checks the current validator and,
// if the next validator is installed, it passes control to it.
// Returns true if all validators are successful.
func (b *baseValidator) Validate() bool {
	if b.next != nil {
		return b.next.Validate()
	}
	return true
}
