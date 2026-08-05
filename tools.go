package slug

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/goloop/t13n/v2/lang"
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

	// A rune with no transliteration (an emoji, a symbol outside the tables)
	// would otherwise vanish and glue its neighbours together. Turn a visible
	// one into a word boundary so it splits words like any other punctuation
	// ("fire🔥sale" -> "fire-sale"); keep invisible format characters
	// (ZWJ/ZWNJ/soft hyphen) glued, matching the "split on non-letters" model.
	if ts.Value == "" {
		if unicode.IsGraphic(ts.Curr) {
			return " ", 0, true
		}
		return "", 0, true
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

// caseMode selects the letter case applied while building the slug, so the
// result is produced in a single pass instead of a second ToLower/ToUpper.
type caseMode int8

const (
	casePreserve caseMode = iota // keep the transliterated case
	caseLower                    // force lower case
	caseUpper                    // force upper case
)

// applyCase returns the ASCII byte b in the requested case. The slug body is
// ASCII, so byte-wise casing is exact.
func applyCase(b byte, mode caseMode) byte {
	switch mode {
	case caseLower:
		if b >= 'A' && b <= 'Z' {
			return b + ('a' - 'A')
		}
	case caseUpper:
		if b >= 'a' && b <= 'z' {
			return b - ('a' - 'A')
		}
	}

	return b
}

// writeCased writes s to b applying the case mode. For casePreserve it is a
// plain WriteString.
func writeCased(b *strings.Builder, s string, mode caseMode) {
	if mode == casePreserve {
		b.WriteString(s)
		return
	}
	for i := 0; i < len(s); i++ {
		b.WriteByte(applyCase(s[i], mode))
	}
}

// render turns the transliterated ASCII string into the final slug in the
// requested case. The common case (no length limit) is a single Builder pass;
// a MaxLength needs word boundaries, so it goes through tokenize + assemble.
func (c config) render(raw string, mode caseMode) string {
	if c.maxLen > 0 {
		return c.assemble(tokenize(raw), mode)
	}
	return c.join(raw, mode)
}

// join keeps letter/digit runs and inserts a single separator between them
// in one pass, without leading, trailing or doubled separators. It returns
// the fallback when nothing survives.
func (c config) join(raw string, mode caseMode) string {
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
		b.WriteByte(applyCase(raw[i], mode))
		pendingSep = false
	}

	if b.Len() == 0 {
		return c.casedFallback(mode)
	}

	return b.String()
}

// casedFallback returns the fallback in the requested case. The fallback is
// stored as Make would render it, so Lower, Upper and MakeURL have to case it
// on the way out; otherwise input that produces nothing escapes the case the
// caller asked for, and a URL slug could come back capitalised.
func (c config) casedFallback(mode caseMode) string {
	switch mode {
	case caseLower:
		return strings.ToLower(c.fallback)
	case caseUpper:
		return strings.ToUpper(c.fallback)
	}

	return c.fallback
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

// assemble joins words with the configured separator in the requested case,
// honouring MaxLength (cutting on a word boundary) and returning the fallback
// for empty input.
func (c config) assemble(tokens []string, mode caseMode) string {
	if len(tokens) == 0 {
		return c.casedFallback(mode)
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
				writeCased(&b, tok[:c.maxLen], mode)
			}
			break
		}

		if i > 0 {
			b.WriteString(c.separator)
		}
		writeCased(&b, tok, mode)
	}

	// Non-empty tokens always yield at least one character (the empty case
	// is handled above), so the result is never empty here.
	return b.String()
}

// withSuffix appends the uniqueness suffix (separator + n) to base and reports
// whether a candidate within MaxLength could be formed. When a MaxLength is
// set the base is shortened whole word by whole word so the suffixed slug
// still fits; if even the bare number would overflow the limit, ok is false
// (no valid candidate exists — the length invariant is never broken).
func (c config) withSuffix(base string, n int) (string, bool) {
	num := strconv.Itoa(n)

	// The bare number is the shortest possible candidate; if it does not fit,
	// nothing longer will either.
	if c.maxLen > 0 && len(num) > c.maxLen {
		return "", false
	}

	if base == "" {
		return num, true
	}

	suffix := c.separator + num
	if c.maxLen > 0 && len(base)+len(suffix) > c.maxLen {
		keep := c.maxLen - len(suffix)
		if keep <= 0 {
			// The suffix alone fills the limit: fall back to the bare number
			// (already checked to fit above).
			return num, true
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

	return base + suffix, true
}
