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

// repoName marks frames that belong to this repository and should be skipped.
var repoName = "github.com/chg1f/errorx"

// frame stores one runtime caller frame before it is rendered for logging.
type frame struct {
	File string
	Line int
	Func string
}

// stack is the default stack implementation exported through errorx.Stack.
type stack []frame

// LogValue renders captured frames as a compact single-line string.
func (s stack) LogValue() slog.Value {
	if len(s) == 0 {
		return slog.StringValue("")
	}
	var b strings.Builder
	for i := range s {
		if i != 0 {
			b.WriteString(" | ")
		}
		name := s[i].Func
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		b.WriteString(filepath.Base(s[i].File))
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(s[i].Line))
		b.WriteByte('@')
		b.WriteString(name)
	}
	return slog.StringValue(b.String())
}

// Stacktrace captures runtime callers while skipping internal repository frames.
func Stacktrace() errorx.Stack {
	frames := make([]frame, 0, Depth)
	pcs := make([]uintptr, Depth)
	count := runtime.Callers(2, pcs)
	iter := runtime.CallersFrames(pcs[:count])
	for len(frames) < Depth {
		item, more := iter.Next()
		if !strings.Contains(item.Function, repoName) {
			frames = append(frames, frame{
				File: item.File,
				Line: item.Line,
				Func: item.Function,
			})
		}
		if !more {
			break
		}
	}
	return stack(frames)
}

// init installs the stack provider into the root errorx package.
func init() {
	errorx.Stacktrace = Stacktrace
}
