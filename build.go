package errorx

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
)

type Stack interface {
	fmt.Stringer
	slog.LogValuer
	json.Marshaler
	// iter.Seq[string]
}

var Stacktrace = func() Stack { return nil }

type Builder[T comparable] Error[T]

func (eb Builder[T]) clone() Builder[T] {
	var nb Builder[T]
	nb.code = eb.code
	nb.values = make(map[string]any, len(eb.values))
	for k := range eb.values {
		nb.values[k] = eb.values[k]
	}
	return nb
}

func Code[T comparable](code T) Builder[T] {
	var eb Builder[T]
	eb.code = code
	eb.values = map[string]any{}
	return eb
}

func (eb Builder[T]) With(key string, value any) Builder[T] {
	n := eb.clone()
	n.values[key] = value
	return n
}

var Unspecified = struct{}{}

func New(text string) error { return Code(Unspecified).New(text) }

func (eb Builder[T]) New(msg string) error {
	ex := eb.clone()
	ex.message = msg
	ex.stack = Stacktrace()
	return (*Error[T])(&ex)
}

func Errorf(format string, a ...any) error { return Code(Unspecified).Errorf(format, a...) }

func (eb Builder[T]) Errorf(format string, a ...any) error {
	return eb.New(fmt.Sprintf(format, a...))
}

func Wrap(err error) error { return Code(Unspecified).Wrap(err) }

func (eb Builder[T]) Wrap(err error) error {
	if err != nil {
		ex := eb.clone()
		ex.wrapped = err
		ex.stack = Stacktrace()
		return (*Error[T])(&ex)
	}
	return nil
}

func Join(errs ...error) error { return Code(Unspecified).Join(errs...) }

func (eb Builder[T]) Join(errs ...error) error {
	return eb.Wrap(errors.Join(errs...))
}
