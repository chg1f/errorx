package errorx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"iter"
	"log/slog"
	"slices"
)

var Length = 1024

type Error[T comparable] struct {
	code    T
	values  map[string]any
	wrapped error
	message string

	stack Stack
}

func (ex *Error[T]) Error() string {
	buf := bytes.NewBuffer(make([]byte, 0, Length))
	fmt.Fprintf(buf, "#%v %s", ex.code, ex.String())
	for k, v := range ex.Items() {
		fmt.Fprintf(buf, " %s=%v", k, v)
	}
	return buf.String()
}

var _ error = &Error[struct{}]{}

func (ex *Error[T]) Unwrap() error { return ex.wrapped }

var _ interface{ Unwrap() error } = &Error[struct{}]{}

func (ex *Error[T]) Is(err error) bool { return ex.wrapped == err }

var _ interface{ Is(error) bool } = &Error[struct{}]{}

func (ex *Error[T]) String() string {
	if ex.wrapped != nil {
		if ex.message != "" {
			return ex.message + "; " + ex.wrapped.Error()
		}
		return ex.wrapped.Error()
	}
	return ex.message
}

var _ fmt.Stringer = &Error[struct{}]{}

func (ex Error[T]) Code() T { return ex.code }

func (ex *Error[T]) Get(key string) (any, bool) {
	v, ok := ex.values[key]
	return v, ok
}

func (ex *Error[T]) Items() iter.Seq2[string, any] {
	return func(yield func(string, any) bool) {
		for k, v := range ex.values {
			if !yield(k, v) {
				return
			}
		}
	}
}

func (ex Error[T]) Stacktrace() Stack {
	return ex.stack
}

func (ex *Error[T]) LogValue() slog.Value {
	attrs := []slog.Attr{
		slog.String("code", fmt.Sprintf("%v", ex.code)),
	}
	if ex.wrapped != nil {
		if ex.message != "" {
			slog.String("message", ex.message)
		}
		attrs = append(attrs, slog.Any("wrapped", ex.wrapped))
	} else {
		slog.String("message", ex.message)
	}
	if ex.stack != nil {
		attrs = append(attrs, slog.Any("stack", ex.stack))
	}
	return slog.GroupValue(attrs...)
}

var _ slog.LogValuer = &Error[struct{}]{}

func (ex *Error[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"code":    ex.code,
		"stack":   ex.stack,
		"values":  ex.values,
		"wrapped": ex.wrapped,
		"message": ex.message,
	})
}

var _ json.Marshaler = &Error[struct{}]{}

func (ex *Error[T]) In(codes ...T) bool {
	return slices.Contains(codes, ex.code)
}

func Be[T comparable](err error) *Error[T] {
	if err == nil {
		return nil
	}
	ex, ok := err.(*Error[T])
	if !ok {
		var empty T
		ex = Code(empty).Wrap(err).(*Error[T])
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

func Get[T comparable](err error, key string) (any, bool) {
	for {
		if err == nil {
			return nil, false
		}
		if x, ok := err.(*Error[T]); ok {
			return x.Get(key)
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if v, ok := Get[T](err, key); ok {
					return v, true
				}
			}
			return nil, false
		default:
			return nil, false
		}
	}
}
