package stacktrace

import (
	"log/slog"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/chg1f/errorx/v2"
)

// Depth bounds the number of frames captured by the default provider.
var Depth = 16

// PackageNames holds the default package filters used by Stacktrace.
var PackageNames = []string{
	"github.com/chg1f/errorx/v2",
}

// stack is the default stack implementation exported through errorx.Stack.
type stack []uintptr

// LogValue renders captured frames as a compact single-line string.
func (s stack) LogValue() slog.Value {
	if len(s) == 0 {
		return slog.StringValue("")
	}
	var b strings.Builder
	frames := s.Frames()
	for i := 0; ; i++ {
		frame, more := frames.Next()
		if frame.File == "" {
			break
		}
		if i != 0 {
			b.WriteString(" | ")
		}
		name := frame.Function
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		b.WriteString(filepath.Base(frame.File))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(frame.Line))
		b.WriteByte('@')
		b.WriteString(name)
		if !more {
			break
		}
	}
	return slog.StringValue(b.String())
}

// Frames rebuilds a fresh runtime.Frames iterator for callers that need frame access.
func (s stack) Frames() runtime.Frames {
	frames := runtime.CallersFrames([]uintptr(s))
	if frames == nil {
		return runtime.Frames{}
	}
	return *frames
}

// Stacktrace builds a stack provider after preprocessing the current PackageNames.
func Stacktrace() func() errorx.Stack {
	return func() errorx.Stack {
		preservedPcs := make([]uintptr, 0, Depth)
		pcs := make([]uintptr, Depth)
		count := runtime.Callers(2, pcs)
		iter := runtime.CallersFrames(pcs[:count])
		for len(preservedPcs) < Depth {
			item, more := iter.Next()
			if !skipFrame(item, PackageNames...) {
				preservedPcs = append(preservedPcs, item.PC)
			}
			if !more {
				break
			}
		}
		return stack(preservedPcs)
	}
}

// skipFrame reports whether the frame belongs to one of the filtered packages.
func skipFrame(frame runtime.Frame, pkgNames ...string) bool {
	for _, pkgName := range pkgNames {
		if strings.Contains(frame.Function, pkgName) {
			return true
		}
	}
	return false
}
