package errorx_test

import (
	"errors"
	"log/slog"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuilderMethods verifies builder chaining preserves code and metadata.
func TestBuilderMethods(t *testing.T) {
	err := errorx.WithCode("invalid").
		New("invalid input", slog.String("field", "email"))

	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	assert.Equal(t, "invalid", ex.Code())
	assert.EqualError(t, err, "#invalid invalid input")
}

// TestWrapAndIs verifies wrapping keeps code and cause discoverable.
func TestWrapAndIs(t *testing.T) {
	base := errors.New("disk failure")
	err := errorx.WithCode("missing").Wrap(base, "load config")

	assert.ErrorIs(t, err, errorx.WithCode("missing").New(""))
	assert.True(t, errorx.In(errorx.Be[string](err), "missing"))
	assert.ErrorIs(t, err, base)
	assert.Nil(t, errorx.Wrap(nil, ""))
	assert.EqualError(t, err, "#missing load config; disk failure")
}

// TestPackageLevelConstructors verifies the package-level convenience functions.
func TestPackageLevelConstructors(t *testing.T) {
	assert.EqualError(t, errorx.New("plain"), "plain")
	err := errorx.WithCode(errorx.Unspecified).New("override", slog.String("k", "v"))
	assert.EqualError(t, err, "override")
}
