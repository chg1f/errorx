package errorx

type unspecified struct{}

var Unspecified = unspecified{}

// isUnspecified hides the sentinel code from Error output.
func isUnspecified(code any) bool {
	_, ok := code.(unspecified)
	return ok
}
