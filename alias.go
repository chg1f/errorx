package errorx

import (
	"errors"
)

var Unspecified = struct{}{}

// alias errors.Unwrap
func Unwrap(err error) error {
	return errors.Unwrap(err)
}

// alias errors.Is
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// alias errors.As
func As(err error, target any) bool {
	return errors.As(err, target)
}

// hijack errors.New
func New(text string) error {
	return Code(Unspecified).New(text)
}

// hijack errors.Wrap
func Join(errs ...error) error {
	return Code(Unspecified).Join(errs...)
}

// hijack fmt.Errorf
func Errorf(format string, a ...any) error {
	return Code(Unspecified).Errorf(format, a...)
}
