package errorx

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"
)

type unspecified struct{}

var Unspecified = unspecified{}

// Length controls the initial buffer size used by Error.Error.
var Length = 256

// Error represents a typed error node that always carries a code.
type Error[T comparable] struct {
	code T

	message string
	cause   error

	attrs []slog.Attr
	stack Stack
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

// Message returns the static message without attrs, code, stack, or wrapped cause.
func (ex *Error[T]) Message() string {
	if ex == nil {
		return ""
	}
	return ex.message
}

// Stack returns the captured stack, if a stack provider was registered.
func (ex *Error[T]) Stack() Stack {
	return ex.stack
}

// Unwrap exposes wrapped causes to the standard errors package.
func (ex *Error[T]) Unwrap() error {
	return ex.cause
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

// In reports whether the stored code matches one of the provided codes.
func (ex *Error[T]) In(codes []T) bool {
	if ex == nil {
		return false
	}
	return slices.Contains(codes, ex.code)
}

// Error formats the stack location, code, rendered body, and wrapped cause.
func (ex *Error[T]) Error() string {
	if ex == nil {
		return "<nil>"
	}
	var buf strings.Builder
	buf.Grow(Length)
	if _, ok := any(ex.code).(unspecified); !ok {
		_, _ = fmt.Fprintf(&buf, "#%v", ex.code)
	}
	buf.WriteString(ex.String())
	if ex.stack != nil {
		buf.WriteString(" [")
		buf.WriteString(ex.stack.String())
		buf.WriteByte(']')
	}
	return buf.String()
}

var _ error = (*Error[struct{}])(nil)

// String formats the message, attrs, and wrapped cause without rendering the code.
func (ex *Error[T]) String() string {
	if ex == nil {
		return "<nil>"
	}
	var buf strings.Builder
	if ex.message != "" {
		buf.WriteString(ex.message)
	}
	var attrs strings.Builder
	for i := range ex.attrs {
		attr := ex.attrs[i]
		attr.Value = attr.Value.Resolve()
		if attr.Key == "" {
			continue
		}
		if attrs.Len() != 0 {
			attrs.WriteByte(' ')
		}
		attrs.WriteString(attr.Key)
		attrs.WriteByte('=')
		_, _ = fmt.Fprint(&attrs, attr.Value.Any())
	}
	if attrs.Len() != 0 {
		buf.WriteByte('(')
		buf.WriteString(attrs.String())
		buf.WriteByte(')')
	}
	if ex.cause != nil {
		if buf.Len() != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ex.cause.Error())
	}
	return buf.String()
}

var _ fmt.Stringer = (*Error[struct{}])(nil)

// LogValue exposes the error as a structured slog group.
func (ex *Error[T]) LogValue() slog.Value {
	if ex == nil {
		return slog.AnyValue(nil)
	}
	attrs := []slog.Attr{
		slog.Any("code", ex.code),
		slog.String("message", ex.message),
	}
	if ex.cause != nil {
		attrs = append(attrs, slog.Any("cause", ex.cause))
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
