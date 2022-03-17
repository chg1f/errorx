package errorx

import (
	"errors"
	"fmt"
	"strings"
)

type WeakError struct {
	i interface{}
}

func (e WeakError) Error() string {
	if ret, ok := e.i.(interface {
		String() string
	}); ok {
		return ret.String()
	}
	return fmt.Sprintf("%v", e.i)
}
func (e WeakError) Interface() interface{} {
	return e.i
}
func (e WeakError) As(interface{}) bool {
	return false
}
func (e WeakError) Is(error) bool {
	return false
}

func AsError(i interface{}) error {
	return WeakError{i: i}
}

type CompressError struct {
	x []error
}

func (e CompressError) Error() string {
	t := make([]string, 0, len(e.x))
	for _, e := range e.x {
		t = append(t, e.Error())
	}
	return strings.Join(t, "; ")
}
func (e CompressError) Append(n error) {
	e.x = append(e.x, n)
}
func (e CompressError) As(t interface{}) bool {
	for _, e := range e.x {
		if errors.As(e, t) {
			return true
		}
	}
	return false
}
func (e CompressError) Is(t error) bool {
	for _, e := range e.x {
		if errors.Is(e, t) {
			return true
		}
	}
	return false
}

func Compress(es ...error) error {
	return CompressError{x: es}
}
