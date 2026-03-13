package errorx_test

import (
	"errors"
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
	err := errorx.WithCode("invalid").New("invalid input")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)

	assert.True(t, errorx.In(ex, "invalid"))
	assert.False(t, errorx.In(ex, "missing"))
	assert.False(t, errorx.In[string](nil, "invalid"))
}
