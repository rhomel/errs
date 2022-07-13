package errs

// A package to provide both constant errors and arbitrarily long chained
// errors of multiple types.

import (
	"errors"
	"reflect"
)

// Const represents a constant error.
type Const string

func (e Const) Error() string {
	return string(e)
}

// Error adds the ability to have a chain of errors that are compatible with
// the standard library errors package Is() and As() functions.
// TODO: maybe include stacktraces
type Error struct {
	Current   error
	Next      error
	formatter Formatter
}

func NewError(current, origin error) *Error {
	return NewErrorF(current, origin, DefaultFormatter)
}

type Formatter func(error, error) string

func DefaultFormatter(current, origin error) string {
	return current.Error() + ": " + origin.Error()
}

func NewErrorF(current, origin error, formatter Formatter) *Error {
	return &Error{
		Current:   current,
		Next:      origin,
		formatter: formatter,
	}
}

func (e *Error) Error() string {
	return e.formatter(e.Current, e.Next)
}

func (e *Error) Unwrap() error {
	return e.Next
}

func (e *Error) Is(target error) bool {
	return errors.Is(e.Current, target)
}

func (e *Error) As(target interface{}) bool {
	// NOTE: this code is taken from the Golang errors package implementation
	// for the As function. It appears there is no other way to dereference a
	// double pointer from an empty interface without the reflect package.
	val := reflect.ValueOf(target)
	typ := val.Type()
	if typ.Kind() != reflect.Ptr || val.IsNil() {
		panic("errors: target must be a non-nil pointer")
	}
	targetType := typ.Elem()
	if targetType.Kind() != reflect.Interface && !targetType.Implements(errorType) {
		panic("errors: *target must be interface or implement error")
	}
	err := e.Current
	if reflect.TypeOf(err).AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(err))
		return true
	}
	return false
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()
