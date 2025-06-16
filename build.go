package errorx

import (
	"errors"
	"fmt"
)

var Stacktrace = func() Stack { return nil }

type Builder[T comparable] Error[T]

func (eb Builder[T]) clone() Builder[T] {
	return Builder[T]{
		code: eb.code,
	}
}

func Code[T comparable](code T) Builder[T] {
	var eb Builder[T]
	eb.code = code
	return eb
}

func (eb Builder[T]) Message(msg string) Builder[T] {
	n := eb.clone()
	n.msg = msg
	return n
}

func (eb Builder[T]) Override(msg string) Builder[T] {
	n := eb.clone()
	n.override = msg
	return n
}

var Unspecified = struct{}{}

func New(text string) error { return Code(Unspecified).New(text) }

func (eb Builder[T]) New(msg string) error {
	ex := eb.clone()
	ex.msg = msg
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
		ex.err = err
		ex.stack = Stacktrace()
		return (*Error[T])(&ex)
	}
	return nil
}

func Join(errs ...error) error { return Code(Unspecified).Join(errs...) }

func (eb Builder[T]) Join(errs ...error) error {
	return eb.Wrap(errors.Join(errs...))
}
