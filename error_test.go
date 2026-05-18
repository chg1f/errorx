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
	assert.EqualError(t, err, "#invalid invalid input(field=email)")
}

// TestWrapAndIs verifies wrapping keeps code and cause discoverable.
func TestWrapAndIs(t *testing.T) {
	base := errors.New("disk failure")
	err := errorx.WithCode("missing").Wrap(base, "load config")

	assert.ErrorIs(t, err, errorx.WithCode("missing").New(""))
	assert.True(t, errorx.In(errorx.Be[string](err), []string{"missing"}))
	assert.ErrorIs(t, err, base)
	assert.Nil(t, errorx.Wrap(nil, ""))
	assert.Contains(t, err.Error(), "#missing load config, disk failure")
}

// TestPackageLevelConstructors verifies the package-level convenience functions.
func TestPackageLevelConstructors(t *testing.T) {
	assert.EqualError(t, errorx.New("plain"), "plain")
	err := errorx.WithCode(errorx.Unspecified).New("override", slog.String("k", "v"))
	assert.EqualError(t, err, "override(k=v)")
}

// TestMessageStringAndError verifies the layered renderers expose the expected detail.
func TestMessageStringAndError(t *testing.T) {
	base := errors.New("disk failure")
	err := errorx.WithCode("missing").
		Wrap(base, "load config", slog.String("file", "app.yaml"))

	ex := errorx.Be[string](err)
	require.NotNil(t, ex)

	assert.Equal(t, "load config", ex.Message())
	assert.Equal(t, "load config(file=app.yaml)", ex.String())
	assert.Equal(t, "#missing load config(file=app.yaml), disk failure", ex.Error())
}

// TestJoin verifies flattened aggregation and stable single-line formatting.
func TestJoin(t *testing.T) {
	errA := errors.New("alpha")
	errB := errorx.WithCode("invalid").New("beta")
	errC := errors.New("gamma")

	err := errorx.Join(nil, errA, errorx.Join(errB, nil, errC))

	require.Error(t, err)
	assert.EqualError(t, err, "alpha; #invalid beta; gamma")
	assert.ErrorIs(t, err, errA)
	assert.ErrorIs(t, err, errB)
	assert.ErrorIs(t, err, errC)

	multi, ok := err.(interface{ Unwrap() []error })
	require.True(t, ok)
	assert.Len(t, multi.Unwrap(), 3)
}

// TestJoinSingle verifies Join avoids wrapping a single remaining error.
func TestJoinSingle(t *testing.T) {
	base := errors.New("only")

	err := errorx.Join(nil, base, nil)

	assert.Same(t, base, err)
}
