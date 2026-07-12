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

	// defaultCase is the case Make (and therefore MakeUnique and
	// TryMakeUnique) applies. The zero value keeps the transliterated case;
	// WithLowercase and WithUppercase force a case. Lower and Upper always
	// force their own case regardless of this setting.
	defaultCase caseMode
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
// concatenates the words without any delimiter.
//
// The separator is not validated. For predictable, URL-safe slugs it should
// contain only URL-safe, non-alphanumeric characters (for example "-" or
// "_"). A separator that contains letters or digits merges into the slug body
// and breaks IsValid; unsafe characters (spaces, ".", "/", "%", ...) produce
// slugs that are not URL-safe.
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
//
// The fallback is itself normalized through the slug pipeline at construction
// time, so it always obeys the canonical-slug and MaxLength invariants. A
// fallback that normalizes to nothing (for example "!!!") is treated as no
// fallback.
func WithFallback(s string) Option {
	return func(c *config) { c.fallback = s }
}

// WithLowercase makes Make (and MakeUnique/TryMakeUnique built on it) return a
// lower-case slug, which is the usual convention for URL slugs. Without it the
// transliterated case is preserved; use Lower for a one-off lower-case slug.
func WithLowercase() Option {
	return func(c *config) { c.defaultCase = caseLower }
}

// WithUppercase makes Make (and MakeUnique/TryMakeUnique built on it) return an
// upper-case slug. Without it the transliterated case is preserved; use Upper
// for a one-off upper-case slug.
func WithUppercase() Option {
	return func(c *config) { c.defaultCase = caseUpper }
}
