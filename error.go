package errorx

import (
	"errors"
	"strings"
)

type CompressError struct {
	Errors []error
}

func (ce *CompressError) Is(t error) bool {
	if t == nil {
		if ce == nil {
			return true
		}
	} else if ce.Errors == nil {
		return false
	}
	for _, e := range ce.Errors {
		if errors.Is(e, t) {
			return true
		}
	}
	return false
}
func (ce *CompressError) As(t interface{}) bool {
	if ce == nil {
		return false
	}
	for _, e := range ce.Errors {
		if errors.As(e, t) {
			return true
		}
	}
	return false
}
func (ce *CompressError) Error() string {
	temp := make([]string, 0, len(ce.Errors))
	for _, e := range ce.Errors {
		if e != nil {
			temp = append(temp, e.Error())
		}
	}
	return strings.Join(temp, "; ")
}

func Compress(es ...error) error {
	ce := CompressError{
		Errors: make([]error, 0, len(es)),
	}
	for _, e := range es {
		if e != nil {
			if t, ok := e.(*CompressError); ok {
				ce.Errors = append(ce.Errors, t.Errors...)
				continue
			}
			ce.Errors = append(ce.Errors, e)
		}
	}
	if len(ce.Errors) > 0 {
		return &ce
	}
	return nil
}

func Shrink(es ...error) error {
	return errors.New(Compress(es...).Error())
}
