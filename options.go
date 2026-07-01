package slug

// DefaultSeparator is the word separator used when no other is configured.
const DefaultSeparator = "-"

// config holds the resolved settings of a Slug. It is populated once,
// in New, from the supplied options and is never mutated afterwards,
// which is what makes a *Slug safe for concurrent use.
type config struct {
	// lang is the transliteration language code (see the lang subpackage).
	// The zero value (lang.None) applies no regional rules.
	lang string

	// separator joins the words of the resulting slug. Defaults to "-".
	// It should not contain letters or digits; such characters would be
	// treated as slug content on the next pass and break IsValid.
	separator string

	// maxLen limits the slug length in bytes (bytes equal runes here, as
	// the output is ASCII). Zero means unlimited. Truncation happens on a
	// word boundary so a word is never split in the middle.
	maxLen int

	// fallback is returned when the input yields an empty slug (input made
	// up entirely of separators/dropped characters). Empty by default.
	fallback string
}

// An Option configures a Slug in New.
type Option func(*config)

// WithLang sets the transliteration language code. Use the constants of
// the lang subpackage, e.g. WithLang(lang.UK). An unknown code falls back
// to language-neutral transliteration (lang.None).
func WithLang(code string) Option {
	return func(c *config) { c.lang = code }
}

// WithSeparator sets the word separator (default "-"). An empty separator
// concatenates the words without any delimiter. The separator should not
// contain letters or digits.
func WithSeparator(sep string) Option {
	return func(c *config) { c.separator = sep }
}

// WithMaxLength limits the slug to at most n bytes, cutting on a word
// boundary so a word is never split. Values <= 0 are ignored (unlimited).
func WithMaxLength(n int) Option {
	return func(c *config) {
		if n > 0 {
			c.maxLen = n
		}
	}
}

// WithFallback sets the value returned when the input produces an empty
// slug (e.g. Make("!!!")). Without it such input yields an empty string.
func WithFallback(s string) Option {
	return func(c *config) { c.fallback = s }
}
