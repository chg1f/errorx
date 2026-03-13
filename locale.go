package errorx

import (
	"log/slog"

	"golang.org/x/text/language"
)

// Localize resolves a localized message using message, attributes, and languages.
var Localize = func(message string, attrs []slog.Attr, langs ...language.Tag) (string, bool) {
	return "", false
}
