package errorx

import (
	"fmt"
	"log/slog"
)

// Stack exposes captured frames.
type Stack interface {
	fmt.Stringer
	slog.LogValuer
}

// Stacktrace is the optional stack provider installed by errorx/stacktrace.
var Stacktrace = func() Stack {
	return nil
}
