package errorx

// Frame stores one captured stack frame.
type Frame struct {
	File string
	Line int
	Func string
}

// Stack exposes captured frames.
type Stack interface {
	Frames() []Frame
}

var Stacktrace = func() Stack {
	return nil
}
