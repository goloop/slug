package slug

import (
	"strconv"
	"strings"

	"github.com/goloop/t13n/lang"
)

// Word replacements for symbols that carry meaning. They are surrounded by
// spaces so the normalization pass wraps them in separators instead of
// gluing them to the neighbouring words (e.g. "r&d" -> "r-and-d").
const (
	symAt    = " at "
	symAnd   = " and "
	symSharp = " sharp "
	symPct   = " pct "
)

// slugRules is the custom transliteration callback for the t13n engine.
// It runs once per input rune, after the base transliteration, and only
// intervenes for a few characters; everything else is passed through
// unchanged and cleaned up later by tokenize.
//
//   - a true apostrophe (one between two letters) is dropped so the letters
//     join: "it's" -> "its", "п'ять" -> "piat";
//   - "@", "&", "#", "%" become their word forms, wrapped in spaces so they
//     end up as standalone words in the slug;
//   - any other character keeps its proposed transliteration; punctuation
//     and separators survive as-is and become word boundaries in tokenize.
func slugRules(ts lang.TransState) (string, int, bool) {
	if ts.IsApostrophe {
		return "", 0, true
	}

	switch ts.Curr {
	case '@':
		return symAt, 0, true
	case '&':
		return symAnd, 0, true
	case '#':
		return symSharp, 0, true
	case '%':
		return symPct, 0, true
	}

	return ts.Value, 0, true
}

// isAlnum reports whether b is an ASCII letter or digit. The transliterated
// text is ASCII, so a byte-wise check is exact and needs no rune decoding.
func isAlnum(b byte) bool {
	return b >= 'a' && b <= 'z' ||
		b >= 'A' && b <= 'Z' ||
		b >= '0' && b <= '9'
}

// render turns the transliterated ASCII string into the final slug. The
// common case (no length limit) is a single Builder pass; a MaxLength needs
// word boundaries, so it goes through tokenize + assemble.
func (c config) render(raw string) string {
	if c.maxLen > 0 {
		return c.assemble(tokenize(raw))
	}
	return c.join(raw)
}

// join keeps letter/digit runs and inserts a single separator between them
// in one pass, without leading, trailing or doubled separators. It returns
// the fallback when nothing survives.
func (c config) join(raw string) string {
	var b strings.Builder
	b.Grow(len(raw))

	pendingSep := false
	for i := 0; i < len(raw); i++ {
		if !isAlnum(raw[i]) {
			pendingSep = true
			continue
		}
		if pendingSep && b.Len() > 0 {
			b.WriteString(c.separator)
		}
		b.WriteByte(raw[i])
		pendingSep = false
	}

	if b.Len() == 0 {
		return c.fallback
	}

	return b.String()
}

// tokenize splits an ASCII string into its maximal runs of letters and
// digits, dropping everything in between. The returned strings share the
// backing array of s (no copying).
func tokenize(s string) []string {
	var tokens []string
	start := -1
	for i := 0; i < len(s); i++ {
		if isAlnum(s[i]) {
			if start < 0 {
				start = i
			}
			continue
		}
		if start >= 0 {
			tokens = append(tokens, s[start:i])
			start = -1
		}
	}
	if start >= 0 {
		tokens = append(tokens, s[start:])
	}

	return tokens
}

// assemble joins words with the configured separator, honouring MaxLength
// (cutting on a word boundary) and returning the fallback for empty input.
func (c config) assemble(tokens []string) string {
	if len(tokens) == 0 {
		return c.fallback
	}

	var b strings.Builder
	total := 0
	for _, t := range tokens {
		total += len(t)
	}
	b.Grow(total + len(c.separator)*(len(tokens)-1))

	for i, tok := range tokens {
		need := len(tok)
		if i > 0 {
			need += len(c.separator)
		}

		if c.maxLen > 0 && b.Len()+need > c.maxLen {
			// The first word alone already exceeds the limit: there is no
			// boundary to cut on, so cut the word itself.
			if i == 0 && len(tok) > c.maxLen {
				b.WriteString(tok[:c.maxLen])
			}
			break
		}

		if i > 0 {
			b.WriteString(c.separator)
		}
		b.WriteString(tok)
	}

	// Non-empty tokens always yield at least one character (the empty case
	// is handled above), so the result is never empty here.
	return b.String()
}

// withSuffix appends the uniqueness suffix (separator + n) to base. When a
// MaxLength is set the base is shortened whole word by whole word so the
// suffixed slug still fits.
func (c config) withSuffix(base string, n int) string {
	num := strconv.Itoa(n)
	if base == "" {
		return num
	}

	suffix := c.separator + num
	if c.maxLen > 0 && len(base)+len(suffix) > c.maxLen {
		keep := c.maxLen - len(suffix)
		if keep <= 0 {
			return num
		}

		// Prefer to drop whole words from the end...
		if c.separator != "" {
			for len(base) > keep {
				idx := strings.LastIndex(base, c.separator)
				if idx < 0 {
					break
				}
				base = base[:idx]
			}
		}

		// ...but if a single word still overflows, shorten the word itself.
		if len(base) > keep {
			base = base[:keep]
		}
	}

	return base + suffix
}
