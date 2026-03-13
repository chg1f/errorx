package errorx

import (
	"fmt"
	"maps"
)

// Builder mirrors Error so building and built values share the same fields.
type Builder[T comparable] Error[T]

// clone performs a defensive copy so builder chains remain immutable.
func (eb *Builder[T]) clone() *Builder[T] {
	var nb Builder[T]
	nb.code = eb.code
	nb.message = eb.message
	nb.locale = eb.locale
	nb.wrapped = append([]error(nil), eb.wrapped...)
	nb.values = make(map[string]any, len(eb.values))
	maps.Copy(nb.values, eb.values)
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
		code:   code,
		values: map[string]any{},
	}
}

// WithMessage replaces the message while preserving other builder fields.
func (eb *Builder[T]) WithMessage(message string) *Builder[T] {
	next := eb.clone()
	next.message = message
	return next
}

// WithMessage creates an unspecified-code builder with the provided message.
func WithMessage(message string) *Builder[unspecified] {
	return WithCode(Unspecified).WithMessage(message)
}

// WithLocale replaces the localization key while preserving other builder fields.
func (eb *Builder[T]) WithLocale(locale string) *Builder[T] {
	next := eb.clone()
	next.locale = locale
	return next
}

// WithLocale creates an unspecified-code builder with the provided localization key.
func WithLocale(locale string) *Builder[unspecified] {
	return WithCode(Unspecified).WithLocale(locale)
}

// WithValues merges values into the builder, overriding duplicated keys.
func (eb *Builder[T]) WithValues(values map[string]any) *Builder[T] {
	next := eb.clone()
	for key, value := range values {
		next.values[key] = value
	}
	return next
}

// WithValues creates an unspecified-code builder with the provided values.
func WithValues(values map[string]any) *Builder[unspecified] {
	return WithCode(Unspecified).WithValues(values)
}

// New finalizes the builder using the provided message.
func (eb *Builder[T]) New(message string) error {
	return eb.WithMessage(message).build()
}

// New creates an unspecified-code error with the provided message.
func New(message string) error {
	return WithCode(Unspecified).New(message)
}

// Errorf finalizes the builder using a formatted message.
func (eb *Builder[T]) Errorf(format string, args ...any) error {
	return eb.New(fmt.Sprintf(format, args...))
}

// Errorf creates an unspecified-code error using fmt.Sprintf formatting.
func Errorf(format string, args ...any) error {
	return WithCode(Unspecified).Errorf(format, args...)
}

// Wrap finalizes the builder while wrapping the provided cause.
func (eb *Builder[T]) Wrap(err error) error {
	if err == nil {
		return nil
	}
	nb := eb.clone()
	nb.wrapped = []error{err}
	return nb.build()
}

// Wrap creates an unspecified-code error that wraps the provided cause.
func Wrap(err error) error {
	return WithCode(Unspecified).Wrap(err)
}

// Joins finalizes the builder while wrapping the joined errors.
func (eb *Builder[T]) Joins(errs ...error) error {
	nb := eb.clone()
	nb.wrapped = compactErrors(errs)
	if len(nb.wrapped) == 0 {
		return nil
	}
	return nb.build()
}

// build finalizes the current builder into an immutable error value.
func (eb *Builder[T]) build() error {
	nb := eb.clone()
	nb.stack = Stacktrace()
	return (*Error[T])(nb)
}

// Joins creates an unspecified-code error wrapping multiple errors.
func Joins(errs ...error) error {
	return WithCode(Unspecified).Joins(errs...)
}
