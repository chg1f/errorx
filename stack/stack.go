package stack

import (
	"runtime"
	"strings"

	"github.com/chg1f/errorx/v2"
)

// Depth bounds the number of frames captured by the default provider.
var Depth = 16

// stack is the default stack implementation exported through errorx.Stack.
type stack struct {
	frames []errorx.Frame
}

// Frames returns a defensive copy of the captured stack frames.
func (s stack) Frames() []errorx.Frame {
	if len(s.frames) == 0 {
		return nil
	}
	frames := make([]errorx.Frame, len(s.frames))
	copy(frames, s.frames)
	return frames
}

// stacktrace records runtime callers while skipping internal errorx frames.
func stacktrace() errorx.Stack {
	frames := make([]errorx.Frame, 0, Depth)
	pcs := make([]uintptr, Depth)
	count := runtime.Callers(2, pcs)
	iter := runtime.CallersFrames(pcs[:count])
	for len(frames) < Depth {
		frame, more := iter.Next()
		if !strings.Contains(frame.Function, "github.com/chg1f/errorx/v2/stack") &&
			!strings.Contains(frame.Function, "github.com/chg1f/errorx/v2.Stacktrace") {
			frames = append(frames, errorx.Frame{
				File: frame.File,
				Line: frame.Line,
				Func: frame.Function,
			})
		}
		if !more {
			break
		}
	}
	return stack{frames: frames}
}

func init() {
	errorx.Stacktrace = stacktrace
}
