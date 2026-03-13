package errorx

import "strings"

type unspecified struct{}

var Unspecified = unspecified{}

// compactErrors removes nil items while preserving order.
func compactErrors(errs []error) []error {
	items := make([]error, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			items = append(items, err)
		}
	}
	return items
}

// joinErrors renders wrapped errors with a semicolon separator.
func joinErrors(errs []error) string {
	parts := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			parts = append(parts, err.Error())
		}
	}
	return strings.Join(parts, ";")
}

// isUnspecified hides the sentinel code from Error output.
func isUnspecified(code any) bool {
	_, ok := code.(unspecified)
	return ok
}
