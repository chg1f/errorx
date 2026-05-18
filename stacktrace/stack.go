package stacktrace

import (
	"log/slog"
	"runtime"
	"strings"

	"github.com/chg1f/errorx/v2"
)

// stack is the default stack implementation exported through errorx.Stack.
type stack []runtime.Frame

// String renders the last captured frame as FileName:Line@FuncName.
func (s stack) String() string {
	if len(s) == 0 {
		return ""
	}
	return Format(s[0])
}

// LogValue renders captured frames as a []string.
func (s stack) LogValue() slog.Value {
	list := make([]string, 0, len(s))
	for _, frame := range s {
		list = append(list, Format(frame))
	}
	return slog.AnyValue(list)
}

// Stacktrace builds a stack provider with the given depth and skip package names.
func Stacktrace(depth int, skipNames ...string) func() errorx.Stack {
	skipNames = append(skipNames, errorx.Package, errorx.PackageName())
	return func() errorx.Stack {
		pcs := make([]uintptr, depth)
		count := runtime.Callers(2, pcs)
		iter := runtime.CallersFrames(pcs[:count])
		preservedFrames := make([]runtime.Frame, 0, depth)
		for len(preservedFrames) < depth {
			frame, more := iter.Next()
			if frame.File != "" {
				skip := false
				for _, skipName := range skipNames {
					if strings.Contains(frame.Function, skipName) {
						skip = true
						break
					}
				}
				if !skip {
					preservedFrames = append(preservedFrames, frame)
				}
			}
			if !more {
				break
			}
		}
		return stack(preservedFrames)
	}
}
