package errorx

import (
	"errors"
	"strings"
)

const (
	ErrorSeparator = "; "
)

type CompressError struct {
	x []error
}

func (e *CompressError) Error() string {
	t := make([]string, 0, len(e.x))
	for i := range e.x {
		t = append(t, e.x[i].Error())
	}
	return strings.Join(t, ErrorSeparator)
}
func (e *CompressError) Is(t error) bool {
	if t == nil {
		return e == t
	}
	if ce := new(CompressError); errors.As(t, &ce) {
		return ce.ForEach(func(_ int, err error) bool {
			if !errors.Is(e, err) {
				return false
			}
			return true
		})
	}
	for i := range e.x {
		if errors.Is(e.x[i], t) {
			return true
		}
	}
	return false
}
func (e *CompressError) As(t interface{}) bool {
	if _, ok := t.(*CompressError); ok {
		t = e
		return true
	}
	for i := range e.x {
		if errors.As(e.x[i], t) {
			return true
		}
	}
	return false
}
func (e *CompressError) Len() int {
	return len(e.x)
}
func (e *CompressError) ForEach(f func(int, error) bool) bool {
	for i := range e.x {
		if !f(i, e.x[i]) {
			return false
		}
	}
	return true
}

func Compress(es ...error) error {
	if len(es) > 0 {
		t := make([]error, 0, len(es))
		for i := range es {
			if es[i] != nil {
				if ce := new(CompressError); errors.As(es[i], &ce) {
					t = append(t, ce.x...)
					continue
				}
				t = append(t, es[i])
			}
		}
		if len(t) > 0 {
			return &CompressError{x: t}
		}
	}
	return nil
}
func Shrink(es ...error) error {
	if err := Compress(es...); err != nil {
		return errors.New(err.Error())
	}
	return nil
}
func ForEach(e error, f func(int, error) bool) bool {
	if ce := new(CompressError); errors.As(e, &ce) {
		return ce.ForEach(f)
	}
	return f(0, e)
}
