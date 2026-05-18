package stacktrace

import (
	"runtime"
	"strconv"
	"strings"
)

// PackageName resolves the caller package import path.
func PackageName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	name := runtime.FuncForPC(pc).Name()
	if name == "" {
		return ""
	}
	lastSlash := strings.LastIndexByte(name, '/')
	start := lastSlash + 1
	dot := strings.IndexByte(name[start:], '.')
	if dot < 0 {
		return ""
	}
	return name[:start+dot]
}

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
