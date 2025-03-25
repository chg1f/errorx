package errorx

var Stacktrace = func() []Frame { return nil }

type Frame struct {
	PC       uintptr
	FileName string
	FileLine int
	FuncName string
}
