package i18n

import (
	"log/slog"
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

	message, ok := localize("invalid", []slog.Attr{slog.String("field", "email")}, language.English)
	require.True(t, ok)
	assert.Equal(t, "email is not allowed", message)

	err := errorx.WithCode("invalid").
		New("invalid", slog.String("field", "email"))

	typed := errorx.Be[string](err)
	require.NotNil(t, typed)
	assert.Equal(t, "email is not allowed", typed.Localize(language.MustParse("zh-CN")))
	assert.Equal(t, "invalid", typed.String())
	assert.EqualError(t, err, "#invalid invalid")
	assert.Error(t, LoadFiles(nil))
}
