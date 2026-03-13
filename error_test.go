package errorx_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/chg1f/errorx/v2/i18n"
	_ "github.com/chg1f/errorx/v2/stack"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// TestBuilderMethods verifies builder chaining preserves code and metadata.
func TestBuilderMethods(t *testing.T) {
	err := errorx.WithCode("invalid").
		New("invalid input", slog.String("field", "email"))

	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	assert.Equal(t, "invalid", ex.Code())
	assert.Equal(t, "invalid input", ex.String())
	assert.Contains(t, err.Error(), "#invalid")
}

// TestWrapAndIs verifies wrapping keeps code and cause discoverable.
func TestWrapAndIs(t *testing.T) {
	base := errors.New("disk failure")
	err := errorx.WithCode("missing").Wrap(base, "load config")

	assert.ErrorIs(t, err, errorx.WithCode("missing").New(""))
	assert.True(t, errorx.In(errorx.Be[string](err), "missing"))
	assert.ErrorIs(t, err, base)
	assert.Nil(t, errorx.Wrap(nil, ""))
	assert.EqualError(t, err, "#missing load config: disk failure")
}

// TestPackageLevelConstructors verifies the package-level convenience functions.
func TestPackageLevelConstructors(t *testing.T) {
	assert.EqualError(t, errorx.New("plain"), "plain")
	err := errorx.WithCode(errorx.Unspecified).New("override", slog.String("k", "v"))
	assert.EqualError(t, err, "override")
}

// TestI18n verifies message resolution is message-key based and locale-aware.
func TestI18n(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/active.zh-CN.json"
	require.NoError(t, os.WriteFile(path, []byte(`[
  {
    "id": "invalid",
    "translation": "字段 {{.field}} 非法"
  }
]`), 0o600))
	require.NoError(t, i18n.LoadFiles(os.DirFS(dir)))

	err := errorx.WithCode("invalid").
		New("invalid", slog.String("field", "email"))

	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	assert.Equal(t, "字段 email 非法", ex.Localize(language.MustParse("zh-CN")))
	assert.Equal(t, "invalid", ex.Localize(language.MustParse("en")))
	assert.Equal(t, "invalid", ex.String())
	assert.EqualError(t, err, "#invalid invalid")
}

// TestStack verifies the stack provider populates stacks by side-effect import.
func TestStack(t *testing.T) {
	err := errorx.WithCode("invalid").New("boom")
	ex := errorx.Be[string](err)
	require.NotNil(t, ex)
	frames := ex.Stack().Frames()
	require.NotEmpty(t, frames)
	assert.NotEmpty(t, frames[0].File)
}
