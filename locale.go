package errorx

import (
	"golang.org/x/text/language"
)

// Localize resolves a localized message using code, language, and metadata.
var Localize = func(lang language.Tag, locale string, values map[string]any) (string, bool) {
	return "", false
}
