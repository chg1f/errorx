package stacktrace

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStack verifies the stack provider populates frames after explicit registration.
func TestStack(t *testing.T) {
	prev := errorx.Stacktrace
	prevPkgNames := append([]string(nil), PackageNames...)
	PackageNames = append([]string(nil), PackageName())
	errorx.Stacktrace = Stacktrace()
	t.Cleanup(func() {
		errorx.Stacktrace = prev
		PackageNames = prevPkgNames
	})

	pkg := PackageName()
	require.NotEmpty(t, pkg)
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	require.NotNil(t, ex.Stack())
	assert.Equal(t, slog.KindString, ex.Stack().LogValue().Kind())
	assert.NotEmpty(t, ex.Stack().LogValue().String())
	st, ok := ex.Stack().(stack)
	require.True(t, ok)
	require.NotEmpty(t, st)

	stackFrames := st.Frames()
	firstFromStack, _ := stackFrames.Next()
	require.NotEmpty(t, firstFromStack.File)
	assert.NotZero(t, firstFromStack.Line)
	assert.False(t, strings.Contains(firstFromStack.Function, pkg))

	frames := ex.Stack().Frames()
	firstFromInterface, _ := frames.Next()
	assert.NotEmpty(t, firstFromInterface.File)
	assert.NotZero(t, firstFromInterface.Line)
	assert.Equal(t, filepath.Base(firstFromStack.File), filepath.Base(firstFromInterface.File))
	assert.Equal(t, firstFromStack.Line, firstFromInterface.Line)
}

// TestErrorRendersStack verifies Error prefixes the first stack frame as File:Line.
func TestErrorRendersStack(t *testing.T) {
	prev := errorx.Stacktrace
	prevPkgNames := append([]string(nil), PackageNames...)
	PackageNames = append([]string(nil), PackageName())
	errorx.Stacktrace = Stacktrace()
	t.Cleanup(func() {
		errorx.Stacktrace = prev
		PackageNames = prevPkgNames
	})

	err := errorx.WithCode("invalid").New("boom")

	assert.Regexp(t, regexp.MustCompile(`^\[[^/\]]+:\d+\]#invalid boom$`), err.Error())
}
