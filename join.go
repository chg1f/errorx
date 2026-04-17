package errorx

import "strings"

// joinError aggregates multiple causes while preserving multi-error traversal.
type joinError struct {
	errs []error
}

// Error joins child error strings with a stable single-line separator.
func (je *joinError) Error() string {
	if je == nil || len(je.errs) == 0 {
		return ""
	}
	var buf strings.Builder
	for i, err := range je.errs {
		if i != 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(err.Error())
	}
	return buf.String()
}

// Unwrap exposes all aggregated causes for errors.Is and errors.As traversal.
func (je *joinError) Unwrap() []error {
	if je == nil {
		return nil
	}
	return je.errs
}

// Join combines non-nil errors and flattens nested multi-error wrappers iteratively.
func Join(errs ...error) error {
	stack := make([]error, 0, len(errs))
	for i := len(errs) - 1; i >= 0; i-- {
		if errs[i] != nil {
			stack = append(stack, errs[i])
		}
	}

	expand := make([]error, 0, len(stack))
	for len(stack) != 0 {
		n := len(stack) - 1
		err := stack[n]
		stack = stack[:n]

		x, ok := err.(interface{ Unwrap() []error })
		if !ok {
			expand = append(expand, err)
			continue
		}

		children := x.Unwrap()
		for i := len(children) - 1; i >= 0; i-- {
			if children[i] != nil {
				stack = append(stack, children[i])
			}
		}
	}

	switch len(expand) {
	case 0:
		return nil
	case 1:
		return expand[0]
	default:
		return &joinError{errs: expand}
	}
}
