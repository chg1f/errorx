package errorx

import (
	"encoding/json"
	"fmt"
	"log/slog"
)

type Stack interface {
	fmt.Stringer
	slog.LogValuer
	json.Marshaler
}

type Error[T comparable] struct {
	code T

	override string
	msg      string
	err      error

	stack Stack
}

func (ex Error[T]) Code() T {
	return ex.code
}

func (ex Error[T]) Stacktrace() Stack {
	return ex.stack
}

func (ex *Error[T]) Error() string {
	if ex.override != "" {
		return ex.override
	}
	if ex.err != nil {
		if ex.msg != "" {
			return fmt.Sprintf("#%v %s;%s", ex.code, ex.msg, ex.err.Error())
		}
		return fmt.Sprintf("#%v %s", ex.code, ex.err.Error())
	}
	return fmt.Sprintf("#%v %s", ex.code, ex.msg)
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

func (ex *Error[T]) String() string {
	return ex.Error()
}

var _ fmt.Stringer = &Error[struct{}]{}

// func (ex *Error[T]) LogValue() slog.Value {
// 	return slog.StringValue(ex.Error())
// }
//
// var _ slog.LogValuer = &Error[struct{}]{}

// func (ex *Error[T]) MarshalJSON() ([]byte, error) {
// 	return []byte(ex.Error()), nil
// }
//
// var _ json.Marshaler = &Error[struct{}]{}

func (ex *Error[T]) In(codes ...T) bool {
	for i := range codes {
		if ex.code == codes[i] {
			return true
		}
	}
	return false
}

func Be[T comparable](err error) *Error[T] {
	if err == nil {
		return nil
	}
	ex, ok := err.(*Error[T])
	if !ok {
		return Code(Unspecified).Wrap(err).(*Error[T])
	}
	return ex
}

func In[T comparable](err error, codes ...T) bool {
	for {
		if err == nil {
			return false
		}
		if x, ok := err.(interface{ In(...T) bool }); ok {
			return x.In(codes...)
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if In(err, codes...) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}
