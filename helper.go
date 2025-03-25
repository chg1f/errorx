package errorx

import "errors"

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

type Comparable[T comparable] interface {
	In(T) bool
}

func In[T comparable](err error, code T) bool {
	for {
		if err == nil {
			return false
		}
		if x, ok := err.(Comparable[T]); ok {
			return x.In(code)
		}
		switch x := err.(type) {
		case interface{ Unwrap() error }:
			err = x.Unwrap()
		case interface{ Unwrap() []error }:
			for _, err := range x.Unwrap() {
				if In(err, code) {
					return true
				}
			}
			return false
		default:
			return false
		}
	}
}

func Be[T comparable](err error) *Error[T] {
	if err == nil {
		return nil
	}
	ex, ok := err.(*Error[T])
	if !ok {
		return build[T]().Wrap(err).(*Error[T])
	}
	return ex
}

func Stack(err error) []Frame {
	if ex := Be[struct{}](err); ex != nil {
		return ex.Stacktrace()
	}
	return []Frame{}
}
