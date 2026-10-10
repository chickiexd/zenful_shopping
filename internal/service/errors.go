package service

import "fmt"

// ValidationError is a problem with the request that the user can fix;
// handlers report it as 400 with its message.
type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string { return e.msg }

func validationErrorf(format string, args ...any) error {
	return &ValidationError{msg: fmt.Sprintf(format, args...)}
}
