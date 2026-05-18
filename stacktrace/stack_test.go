package stacktrace

import (
	"log/slog"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStack verifies the stack provider populates frames after explicit registration.
func TestStack(t *testing.T) {
	prev := errorx.Stacktrace
	pkg := PackageName()
	errorx.Stacktrace = Stacktrace(16, pkg)
	t.Cleanup(func() {
		errorx.Stacktrace = prev
	})

	require.NotEmpty(t, pkg)
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	require.NotNil(t, ex.Stack())
	assert.Equal(t, slog.KindAny, ex.Stack().LogValue().Kind())
	st, ok := ex.Stack().(stack)
	require.True(t, ok)
	require.NotEmpty(t, st)

	firstFromStack := runtime.Frame(st[0])
	require.NotEmpty(t, firstFromStack.File)
	assert.NotZero(t, firstFromStack.Line)
	assert.False(t, strings.Contains(firstFromStack.Function, pkg))
	frames := make([]string, 0, len(st))
	for _, frame := range st {
		if frame.File == "" {
			continue
		}
		frames = append(frames, Format(frame))
	}
	require.NotEmpty(t, frames)
	assert.Equal(t, frames, ex.Stack().LogValue().Any().([]string))
	assert.Equal(t, frames[len(frames)-1], ex.Stack().String())
	assert.Regexp(t, regexp.MustCompile(`^[^/\]]+:\d+@.+$`), ex.Stack().String())
	assert.Contains(t, frames[0], filepath.Base(firstFromStack.File))
	assert.Equal(t, frames, st.LogValue().Any().([]string))
}

// TestErrorRendersStack verifies Error prefixes the rendered stack string.
func TestErrorRendersStack(t *testing.T) {
	prev := errorx.Stacktrace
	errorx.Stacktrace = Stacktrace(16, PackageName())
	t.Cleanup(func() {
		errorx.Stacktrace = prev
	})

	err := errorx.WithCode("invalid").New("boom")

	assert.Regexp(t, regexp.MustCompile(`^\[[^/\]]+:\d+@.+\]#invalidboom$`), err.Error())
}
