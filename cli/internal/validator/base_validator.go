package validator

type BaseValidator struct {
	next Validatable
}

func (b *BaseValidator) SetNext(next Validatable) {
	b.next = next
}

func (b *BaseValidator) Validate() bool {
	if b.next != nil {
		return b.next.Validate()
	}
	return true
}
