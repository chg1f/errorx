package errorx_test

import (
	"errors"
	"fmt"
	"log/slog"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBe verifies typed extraction and default wrapping behavior.
func TestBe(t *testing.T) {
	err := errorx.WithCode("invalid").
		New("invalid input", slog.String("field", "email"))

	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	assert.Equal(t, "invalid", ex.Code())

	defaulted := errorx.Be(errors.New("plain"), errorx.Empty("missing"))
	require.NotNil(t, defaulted)
	assert.Equal(t, "missing", defaulted.Code())

	assert.Nil(t, errorx.Be[string](nil))

	var nilError *errorx.Error[string]
	assert.Equal(t, "missing", nilError.Code(errorx.NaN("missing")))
}

// TestIn verifies code membership checks on typed errors.
func TestIn(t *testing.T) {
	base := errors.New("disk failure")
	err := errorx.WithCode("invalid").Wrap(base, "invalid input")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)

	assert.True(t, ex.In([]string{"invalid"}))
	assert.False(t, ex.In([]string{"missing"}))

	var nilError *errorx.Error[string]
	assert.False(t, nilError.In([]string{"invalid"}))

	assert.True(t, errorx.In(ex, []string{"invalid"}))
	assert.False(t, errorx.In(ex, []string{"missing"}))
	assert.False(t, errorx.In(nil, []string{"invalid"}))
	assert.True(t, errorx.In(err, []string{"invalid"}))
	assert.True(t, errorx.In(errorx.Join(errors.New("plain"), err), []string{"invalid"}))
	assert.True(t, errorx.In(fmt.Errorf("outer: %w", err), []string{"invalid"}))
}
