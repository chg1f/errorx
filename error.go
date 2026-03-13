package errorx

import (
	"bytes"
	"fmt"
	"log/slog"

	"golang.org/x/text/language"
)

// Length controls the initial buffer size used by Error.Error.
var Length = 256

// Error represents a typed error node that always carries a code.
type Error[T comparable] struct {
	code  T
	attrs []slog.Attr

	message string
	wrapped error

	stack Stack
}

// Error formats the code, resolved message, metadata, and wrapped cause.
func (ex *Error[T]) Error() string {
	if ex == nil {
		return "<nil>"
	}
	buf := bytes.NewBuffer(make([]byte, 0, Length))
	if !isUnspecified(ex.code) {
		fmt.Fprintf(buf, "#%v", ex.code)
	}
	if text := ex.String(); text != "" {
		if buf.Len() > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(buf, "%s", text)
	}
	return buf.String()
}

// String returns the human-readable message chain without metadata decoration.
func (ex *Error[T]) String() string {
	if ex == nil {
		return ""
	}
	switch {
	case ex.message != "" && ex.wrapped != nil:
		return ex.message + ": " + ex.wrapped.Error()
	case ex.message != "":
		return ex.message
	case ex.wrapped != nil:
		return ex.wrapped.Error()
	default:
		return ""
	}
}

// Unwrap exposes wrapped causes to the standard errors package.
func (ex *Error[T]) Unwrap() error {
	return ex.wrapped
}

// Is compares typed error codes to support errors.Is against sentinel builders.
func (ex *Error[T]) Is(target error) bool {
	if ex == nil || target == nil {
		return false
	}
	if x, ok := target.(*Error[T]); ok {
		return ex.code == x.code
	}
	return false
}

// Code returns the stored code or the configured NaN fallback for nil receivers.
func (ex *Error[T]) Code(opts ...func(*CodeOption[T])) T {
	if ex == nil {
		var opt CodeOption[T]
		for _, apply := range opts {
			apply(&opt)
		}
		return opt.NaN
	}
	return ex.code
}

// Stack returns the captured stack, if a stack provider was registered.
func (ex *Error[T]) Stack() Stack {
	return ex.stack
}

// Localize resolves the message for the provided language tags in order.
func (ex *Error[T]) Localize(langs ...language.Tag) string {
	if ex == nil {
		return ""
	}
	s, ok := Localize(ex.message, ex.attrs, langs...)
	if !ok {
		return ex.message
	}
	return s
}
