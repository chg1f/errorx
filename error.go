package errorx

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type Error[T comparable] struct {
	err error
	msg string

	code T

	stack []Frame
}

func (ex Error[T]) Code() T {
	return ex.code
}

func (ex Error[T]) Stacktrace() []Frame {
	return ex.stack
}

func (ex *Error[T]) Error() string {
	if ex.msg != "" {
		return fmt.Sprintf("#%v %s", ex.code, ex.msg)
	}
	return fmt.Sprintf("#%v %s", ex.code, ex.err.Error())
}

var _ error = &Error[struct{}]{}

func (ex *Error[T]) Unwrap() error {
	return ex.err
}

var _ interface{ Unwrap() error } = &Error[struct{}]{}

func (ex *Error[T]) Is(err error) bool {
	return ex.err == err
}

var _ interface{ Is(error) bool } = &Error[struct{}]{}

func (ex *Error[T]) In(code T) bool {
	return ex.code == code
}

var _ Comparable[struct{}] = &Error[struct{}]{}

func (ex *Error[T]) String() string {
	return ex.Error()
}

var _ fmt.Stringer = &Error[struct{}]{}

func (ex *Error[T]) MarshalJSON() ([]byte, error) {
	return []byte(ex.Error()), nil
}

var _ json.Marshaler = &Error[struct{}]{}

func (ex *Error[T]) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Any("code", ex.code),
		slog.String("msg", ex.Error()),
	)
}

var _ slog.LogValuer = &Error[struct{}]{}
