package error

import "fmt"

type Error struct {
	msg string
}

func NewError(format string, a ...any) *Error {
	return &Error{msg: fmt.Sprintf(format, a...)}
}

var _ error = (*Error)(nil)

func (e *Error) Error() string {
	return e.msg
}
