package errorx

import (
	"log/slog"
	"runtime"
)

// Stack exposes captured frames.
type Stack interface {
	slog.LogValuer
	Frames() runtime.Frames
}

// Stacktrace is the optional stack provider installed by errorx/stacktrace.
var Stacktrace = func() Stack {
	return nil
}
