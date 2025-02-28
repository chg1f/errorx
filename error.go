package errorx

import (
	"fmt"
)

type Error[T comparable] Builder[T]

func (ex *Error[T]) Error() string {
	if ex.msg != "" {
		return fmt.Sprintf("#%v %s", ex.code, ex.msg)
	}
	return fmt.Sprintf("#%v %s", ex.code, ex.err.Error())
}

var _ error = &Error[struct{}]{}

func (ex *Error[T]) Unwrap() error { return ex.err }

var _ interface{ Unwrap() error } = &Error[struct{}]{}

func (ex *Error[T]) Is(err error) bool { return ex.err == err }

var _ interface{ Is(error) bool } = &Error[struct{}]{}

func (ex *Error[T]) In(code T, from ...string) bool {
	return ex.code == code && (len(from) == 0 || ex.from == from[0])
}

func (ex Error[T]) From() string { return ex.from }

func (ex Error[T]) Code() T { return ex.code }

func (ex Error[T]) Stack() []Frame { return ex.stack }
