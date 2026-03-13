package errorx

import (
	"log/slog"
)

// Stack exposes captured frames.
type Stack interface {
	slog.LogValuer
}

// Stacktrace is the optional stack provider installed by errorx/stacktrace.
var Stacktrace = func() Stack {
	return nil
}
