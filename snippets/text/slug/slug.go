package slug

import (
	"strings"
	"unicode"
)

// Make returns a URL/path-friendly ASCII slug.
//
// Behavior:
//   - lowercase
//   - converts whitespace/underscore/dash and most punctuation to a single '-'
//   - drops non-ASCII letters/digits (keeps only a-z, 0-9)
//   - trims leading/trailing '-'
//
// This is intentionally simple and dependency-free. If you need transliteration
// (e.g., "你好" -> "ni-hao"), use a dedicated i18n library.
func Make(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	b.Grow(len(s))
	lastDash := false

	for _, r := range s {
		// Keep ASCII letters/digits.
		if r <= unicode.MaxASCII {
			if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
				b.WriteRune(r)
				lastDash = false
				continue
			}
		}

		// Treat separators/punctuation as dash.
		if unicode.IsSpace(r) || r == '-' || r == '_' || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			if b.Len() > 0 && !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
			continue
		}

		// Drop everything else (including non-ASCII letters).
	}

	out := b.String()
	out = strings.Trim(out, "-")
	return out
}
