package stacktrace

import (
	"log/slog"
	"strings"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStack verifies the stack provider populates frames through side-effect import.
func TestStack(t *testing.T) {
	require.NotEmpty(t, repoName)
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	require.NotNil(t, ex.Stack())
	assert.Equal(t, slog.KindString, ex.Stack().LogValue().Kind())
	assert.NotEmpty(t, ex.Stack().LogValue().String())
	st, ok := ex.Stack().(stack)
	require.True(t, ok)
	require.NotEmpty(t, st)
	assert.NotEmpty(t, st[0].File)
	for i := range st {
		assert.False(t, strings.Contains(st[i].Func, repoName))
	}
}
