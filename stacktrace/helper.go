package stacktrace

import (
	"runtime"
	"strconv"
	"strings"
)

// Format formats a frame as FileName:Line@FuncName.
func Format(frame runtime.Frame) string {
	s := frame.File + ":" + strconv.Itoa(frame.Line)
	if name := frame.Function; name != "" {
		if idx := strings.LastIndex(name, "/"); idx >= 0 {
			name = name[idx+1:]
		}
		s += "@" + name
	}
	return s
}
