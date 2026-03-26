package errorx

import (
	"fmt"
	"log/slog"
	"strings"
)

type unspecified struct{}

var Unspecified = unspecified{}

// isUnspecified hides the sentinel code from Error output.
func isUnspecified(code any) bool {
	_, ok := code.(unspecified)
	return ok
}

// formatAttrs renders slog attributes as a compact k=v sequence.
func formatAttrs(attrs []slog.Attr) string {
	var b strings.Builder
	for i := range attrs {
		attr := attrs[i]
		attr.Value = attr.Value.Resolve()
		if attr.Key == "" {
			continue
		}
		if b.Len() != 0 {
			b.WriteByte(' ')
		}
		b.WriteString(attr.Key)
		b.WriteByte('=')
		b.WriteString(fmt.Sprint(attr.Value.Any()))
	}
	return b.String()
}
