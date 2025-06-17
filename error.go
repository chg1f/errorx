package errorx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

var Length = 1024

type Stack interface {
	fmt.Stringer
	slog.LogValuer
	json.Marshaler
}

type Error[T comparable] struct {
	code T

	err      error
	msg      string
	override string

	values map[string]any

	stack Stack
}

func (ex Error[T]) Code() T {
	return ex.code
}

func (ex *Error[T]) Get(key string) (any, bool) {
	if ex.values == nil {
		return nil, false
	}
	v, ok := ex.values[key]
	return v, ok
}

func (ex Error[T]) Stacktrace() Stack {
	return ex.stack
}

func (ex *Error[T]) Error() string {
	bs := make([]byte, 0, Length)
	buf := bytes.NewBuffer(bs)

	fmt.Fprintf(buf, "#%v", ex.code)

	if ex.override != "" {
		buf.WriteString(" " + ex.override)
	} else {
		if ex.err != nil {
			if ex.msg != "" {
				fmt.Fprintf(buf, " %s;%s", ex.msg, ex.err.Error())
			} else {
				fmt.Fprintf(buf, " %s", ex.err.Error())
			}
		} else {
			fmt.Fprintf(buf, " %s", ex.msg)
		}
	}

	if ex.values != nil {
		values := make([]string, 0, len(ex.values))
		for k := range ex.values {
			values = append(values, fmt.Sprintf("%s=%v", k, ex.values[k]))
		}
		buf.WriteString(strings.Join(values, " "))
	}

	return buf.String()
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
	return slices.Contains(codes, ex.code)
	// for i := range codes {
	// 	if ex.code == codes[i] {
	// 		return true
	// 	}
	// }
	// return false
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
