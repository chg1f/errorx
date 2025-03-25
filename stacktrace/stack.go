package stacktrace

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/chg1f/errorx"
)

var (
	Depth = 10

	packageName = reflect.TypeOf(errorx.Error[struct{}]{}).PkgPath()
)

func init() {
	errorx.Stacktrace = stacktrace
}

func stacktrace() []errorx.Frame {
	frames := make([]errorx.Frame, 0, Depth)
	for i := 0; i < Depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, packageName) {
			continue
		}
		frames = append(frames, errorx.Frame{
			PC:       pc,
			FileName: file,
			FileLine: line,
			FuncName: runtime.FuncForPC(pc).Name(),
		})
	}
	return frames
}
