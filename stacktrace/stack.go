package stacktrace

import (
	"fmt"
	"log/slog"
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

type frame struct {
	pc       uintptr
	fileName string
	fileLine int
	funcName string
}

func (f frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.fileName, f.fileLine, f.funcName)
}

type stack []frame

var _ errorx.Stack = stack{}

func (s stack) String() string {
	return fmt.Sprintf("%v", []frame(s))
}

func (s stack) LogValue() slog.Value {
	return slog.AnyValue([]frame(s))
}

func stacktrace() errorx.Stack {
	frames := make([]frame, 0, Depth)
	for i := 0; i < Depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, packageName) {
			continue
		}
		frames = append(frames, frame{
			pc:       pc,
			fileName: file,
			fileLine: line,
			funcName: runtime.FuncForPC(pc).Name(),
		})
	}
	return stack(frames)
}
