package stack_test

import (
	"testing"

	"github.com/chg1f/errorx/v2"
	_ "github.com/chg1f/errorx/v2/stack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStack verifies the stack provider populates frames through side-effect import.
func TestStack(t *testing.T) {
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	frames := ex.Stack().Frames()
	require.NotEmpty(t, frames)
	assert.NotEmpty(t, frames[0].File)
}
