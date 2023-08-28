package errs

import (
	"fmt"
)

// New new an error
func New(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{Message: err.Error()}
}

// Newf new an error with format
func Newf(format string, args ...interface{}) *Error {
	return New(fmt.Errorf(format, args...))
}

// Prefix add prefix to err's router
func (o *Error) Prefix(prefix string) *Error {
	o.Router = prefix + o.Router
	return o
}

// Error error
type Error struct {
	Router  string
	Message string
}

// Error error interface
func (o *Error) Error() string {
	return fmt.Sprintf("%s: %s", o.Router, o.Message)
}
