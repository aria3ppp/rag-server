package error

type ValidationError struct {
	internal error
}

func NewValidationError(internal error) *ValidationError {
	return &ValidationError{internal: internal}
}

var _ error = (*ValidationError)(nil)

func (e *ValidationError) Error() string {
	return e.internal.Error()
}

func (e *ValidationError) Unwrap() error {
	return e.internal
}
