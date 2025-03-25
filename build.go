package errorx

import (
	"errors"
	"fmt"
)

type Builder[T comparable] Error[T]

func build[T comparable]() Builder[T] {
	return Builder[T]{}
}

func (eb Builder[T]) clone() Builder[T] {
	return Builder[T]{
		code: eb.code,
	}
}

func (eb Builder[T]) New(msg string) error {
	ex := eb.clone()
	ex.msg = msg
	ex.stack = stack()
	return (*Error[T])(&ex)
}

func (eb Builder[T]) Errorf(format string, args ...interface{}) error {
	return eb.New(fmt.Sprintf(format, args...))
}

func Wrap(err error) error {
	return Code(Unspecified).Wrap(err)
}

func (eb Builder[T]) Wrap(err error) error {
	if err != nil {
		ex := eb.clone()
		ex.err = err
		ex.stack = stack()
		return (*Error[T])(&ex)
	}
	return nil
}

func (eb Builder[T]) Join(e ...error) error {
	return eb.Wrap(errors.Join(e...))
}

func Code[T comparable](code T) Builder[T] {
	return build[T]().Code(code)
}

func (eb Builder[T]) Code(code T) Builder[T] {
	nb := eb.clone()
	nb.code = code
	return nb
}
