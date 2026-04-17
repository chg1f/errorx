package errorx

import (
	"fmt"
	"log/slog"
	"path/filepath"
	"strconv"
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
	wrapped error

	attrs []slog.Attr
	stack Stack
}

// Message returns the static message without attrs, code, stack, or wrapped cause.
func (ex *Error[T]) Message() string {
	if ex == nil {
		return ""
	}
	return ex.message
}

// String formats the message, attrs, and wrapped cause without rendering the code.
func (ex *Error[T]) String() string {
	if ex == nil {
		return "<nil>"
	}
	return ex.renderBody()
}

// Error formats the stack location, code, rendered body, and wrapped cause.
func (ex *Error[T]) Error() string {
	if ex == nil {
		return "<nil>"
	}
	var buf strings.Builder
	buf.Grow(Length)
	if stack := ex.stackString(); stack != "" {
		buf.WriteByte('[')
		buf.WriteString(stack)
		buf.WriteByte(']')
	}
	if !isUnspecified(ex.code) {
		_, _ = fmt.Fprintf(&buf, "#%v", ex.code)
	}
	if text := ex.renderBody(); text != "" {
		if buf.Len() != 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(text)
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

// renderBody formats message, attrs, and wrapped cause in a stable single-line form.
func (ex *Error[T]) renderBody() string {
	var buf strings.Builder
	if ex.message != "" {
		buf.WriteString(ex.message)
	}
	if attrs := formatAttrs(ex.attrs); attrs != "" {
		buf.WriteByte('(')
		buf.WriteString(attrs)
		buf.WriteByte(')')
	}
	if ex.wrapped != nil {
		if buf.Len() != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(ex.wrapped.Error())
	}
	return buf.String()
}

// stackString renders the first stack frame as File:Line for compact error output.
func (ex *Error[T]) stackString() string {
	if ex == nil || ex.stack == nil {
		return ""
	}
	frames := ex.stack.Frames()
	frame, ok := frames.Next()
	if !ok || frame.File == "" || frame.Line == 0 {
		return ""
	}
	return filepath.Base(frame.File) + ":" + strconv.Itoa(frame.Line)
}
