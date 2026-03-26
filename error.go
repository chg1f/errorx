package errorx

import (
	"fmt"
	"log/slog"
	"strings"
)

// Length controls the initial buffer size used by Error.Error.
var Length = 256

// Error represents a typed error node that always carries a code.
type Error[T comparable] struct {
	code T

	message string
	wrapped error

	attrs []slog.Attr
	stack Stack
}

// Error formats the code, resolved message, metadata, and wrapped cause.
func (ex *Error[T]) Error() string {
	if ex == nil {
		return "<nil>"
	}
	var buf strings.Builder
	buf.Grow(Length)
	if !isUnspecified(ex.code) {
		_, _ = fmt.Fprintf(&buf, "#%v", ex.code)
	}
	if text := ex.message; text != "" {
		if buf.Len() != 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(text)
	}
	if ex.wrapped != nil {
		if buf.Len() != 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(ex.wrapped.Error())
	}
	return buf.String()
}

var _ error = (*Error[struct{}])(nil)

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

// LogValue exposes the error as a structured slog group.
func (ex *Error[T]) LogValue() slog.Value {
	if ex == nil {
		return slog.AnyValue(nil)
	}
	attrs := []slog.Attr{
		slog.Any("code", ex.code),
		slog.String("message", ex.message),
	}
	if ex.wrapped != nil {
		attrs = append(attrs, slog.Any("cause", ex.wrapped))
	}
	if len(ex.attrs) != 0 {
		attrs = append(attrs, slog.Any("attrs", ex.attrs))
	}
	if ex.stack != nil {
		attrs = append(attrs, slog.Any("stack", ex.stack))
	}
	return slog.GroupValue(attrs...)
}

var _ slog.LogValuer = (*Error[struct{}])(nil)
