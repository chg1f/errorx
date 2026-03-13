package i18n

import (
	"os"
	"testing"

	"github.com/chg1f/errorx/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
)

// TestLoadFile verifies bundle loading and error localization behavior.
func TestLoadFile(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(dir+"/active.en.json", []byte(`[ { "id": "invalid", "translation": "{{.field}} is not allowed" } ]`), 0o600))

	require.NoError(t, LoadFiles(os.DirFS(dir)))

	message, ok := localize(language.English, "invalid", map[string]any{
		"field": "email",
	})
	require.True(t, ok)
	assert.Equal(t, "email is not allowed", message)

	err := errorx.WithCode("invalid").
		WithLocale("invalid").
		WithValues(map[string]any{"field": "email"}).
		New("email is invalid")

	typed := errorx.Be[string](err)
	require.NotNil(t, typed)
	assert.Equal(t, "email is not allowed", typed.Localize(language.MustParse("zh-CN")))
	assert.Equal(t, "email is invalid", typed.String())
	assert.EqualError(t, err, "#invalid email is invalid")
	assert.Error(t, LoadFiles(nil))
}
