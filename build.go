package errorx

import "log/slog"

// Builder accumulates error fields before creating an immutable Error value.
type Builder[T comparable] Error[T]

// clone performs a defensive copy so builder chains remain immutable.
func (eb *Builder[T]) clone() *Builder[T] {
	var nb Builder[T]
	nb.code = eb.code
	nb.message = eb.message
	nb.wrapped = eb.wrapped
	nb.attrs = append([]slog.Attr(nil), eb.attrs...)
	nb.stack = eb.stack
	return &nb
}

// WithCode replaces the code while preserving current builder fields.
func (eb *Builder[T]) WithCode(code T) *Builder[T] {
	nb := eb.clone()
	nb.code = code
	return nb
}

// WithCode creates a new typed builder with the provided code.
func WithCode[T comparable](code T) *Builder[T] {
	return &Builder[T]{
		code: code,
	}
}

// WithAttrs appends attributes to the builder in call order.
func (eb *Builder[T]) WithAttrs(attrs ...slog.Attr) *Builder[T] {
	next := eb.clone()
	next.attrs = append(next.attrs, attrs...)
	return next
}

// WithAttrs creates an unspecified-code builder with the provided attributes.
func WithAttrs(attrs ...slog.Attr) *Builder[unspecified] {
	return WithCode(Unspecified).WithAttrs(attrs...)
}

// build finalizes the current builder into an immutable error value.
func (eb *Builder[T]) build() error {
	nb := eb.clone()
	nb.stack = Stacktrace()
	return (*Error[T])(nb)
}

// complete finalizes the builder after applying the last message, cause, and attrs.
func (eb *Builder[T]) complete(message string, wrapped error, attrs ...slog.Attr) error {
	nb := eb.clone()
	nb.message = message
	nb.wrapped = wrapped
	nb.attrs = append(nb.attrs, attrs...)
	return nb.build()
}

// Wrap finalizes the builder while wrapping the provided cause.
func (eb *Builder[T]) Wrap(wrapped error, message string, attrs ...slog.Attr) error {
	if wrapped == nil {
		return nil
	}
	return eb.complete(message, wrapped, attrs...)
}

// Wrap creates an unspecified-code error that wraps the provided cause.
func Wrap(wrapped error, message string, attrs ...slog.Attr) error {
	return WithCode(Unspecified).Wrap(wrapped, message, attrs...)
}

// New finalizes the builder using the provided message.
func (eb *Builder[T]) New(message string, attrs ...slog.Attr) error {
	return eb.complete(message, nil, attrs...)
}

// New creates an unspecified-code error with the provided message.
func New(message string, attrs ...slog.Attr) error {
	return WithCode(Unspecified).New(message, attrs...)
}
