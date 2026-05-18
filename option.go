package errorx

import (
	"errors"
)

// CodeOption defines fallback codes used by Be and Error.Code.
type CodeOption[T comparable] struct {
	NaN   T
	Empty T
}

// Empty configures the code used when wrapping a non-errorx error.
func Empty[T comparable](code T) func(*CodeOption[T]) {
	return func(opt *CodeOption[T]) {
		opt.Empty = code
	}
}

// NaN configures the code returned by Error.Code on a nil receiver.
func NaN[T comparable](code T) func(*CodeOption[T]) {
	return func(opt *CodeOption[T]) {
		opt.NaN = code
	}
}

// Be returns the first typed Error found in the unwrap chain.
// When no typed error exists and a default code is provided, it wraps err with that code.
func Be[T comparable](err error, opts ...func(*CodeOption[T])) *Error[T] {
	if err == nil {
		return nil
	}
	var opt CodeOption[T]
	for _, apply := range opts {
		apply(&opt)
	}
	var ex *Error[T]
	if errors.As(err, &ex) {
		return ex
	}
	if len(opts) == 0 {
		return nil
	}
	ex, _ = WithCode(opt.Empty).Wrap(err, "").(*Error[T])
	return ex
}

// In reports whether any error in the unwrap tree matches one of the provided codes.
func In[T comparable](err error, codes []T, opts ...func(*CodeOption[T])) bool {
	stack := []error{err}
	for len(stack) != 0 {
		n := len(stack) - 1
		node := stack[n]
		stack = stack[:n]
		if node == nil {
			continue
		}
		if ex := Be(node, opts...); ex.In(codes) {
			return true
		}
		if x, ok := node.(interface{ Unwrap() []error }); ok {
			children := x.Unwrap()
			for i := len(children) - 1; i >= 0; i-- {
				if children[i] != nil {
					stack = append(stack, children[i])
				}
			}
			continue
		}
		if x, ok := node.(interface{ Unwrap() error }); ok && x.Unwrap() != nil {
			stack = append(stack, x.Unwrap())
		}
	}
	return false
}
