package stacktrace

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"

	"github.com/chg1f/errorx"
)

var (
	Depth = 10

	pkgPath = reflect.TypeOf(errorx.Error[struct{}]{}).PkgPath()
)

func init() {
	errorx.Stacktrace = Stacktrace
}

type Frame struct {
	pc       uintptr
	fileName string
	fileLine int
	funcName string
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.fileName, f.fileLine, f.funcName)
}

type Stack []Frame

var _ errorx.Stack = Stack{}

func (s Stack) String() string {
	return fmt.Sprintf("%v", []Frame(s))
}

func (s Stack) LogValue() slog.Value {
	return slog.AnyValue([]Frame(s))
}

func (s Stack) MarshalJSON() ([]byte, error) {
	return json.Marshal([]Frame(s))
}

func Stacktrace() errorx.Stack {
	frames := make([]Frame, 0, Depth)
	for i := 0; i < Depth; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		if strings.Contains(file, pkgPath) {
			continue
		}
		frames = append(frames, Frame{
			pc:       pc,
			fileName: file,
			fileLine: line,
			funcName: runtime.FuncForPC(pc).Name(),
		})
	}
	return Stack(frames)
}
