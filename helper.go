package errorx

import (
	"runtime"
	"strings"
)

var Package = PackageName()

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
