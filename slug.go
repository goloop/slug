package slug

import (
	"strings"

	"github.com/goloop/slug/v2/lang"
	"github.com/goloop/t13n/v2"
)

// Slug generates URL-friendly slugs according to its configuration.
//
// A Slug is created with New and is immutable afterwards: all settings are
// resolved from the options at construction time. Because nothing is
// mutated during Make/Lower/Upper/IsValid/MakeUnique, a single *Slug is
// safe for concurrent use by multiple goroutines.
type Slug struct {
	cfg config
}

// New returns a Slug configured by the given options. With no options it
// uses language-neutral transliteration, "-" as the separator, no length
// limit and an empty fallback.
//
//	s := slug.New(
//		slug.WithLang(lang.UK),
//		slug.WithMaxLength(60),
//		slug.WithFallback("post"),
//	)
//	s.Make("Привіт, світ!") // "Pryvit-svit"
func New(opts ...Option) *Slug {
	cfg := config{lang: lang.None, separator: DefaultSeparator}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Slug{cfg: cfg}
}

// Make returns the slug of t. The letter case of the input is preserved;
// use Lower or Upper to force a case.
//
// The input is transliterated to ASCII and then split into words on any
// character that is not a letter or a digit; the words are joined with the
// configured separator. Leading, trailing and repeated separators are
// removed. If the result is empty the configured fallback is returned
// (empty by default).
func (s *Slug) Make(t string) string {
	raw := t13n.Render(s.cfg.lang, t, slugRules)
	return s.cfg.render(raw)
}

// Lower is Make with the result lowercased.
func (s *Slug) Lower(t string) string {
	return strings.ToLower(s.Make(t))
}

// Upper is Make with the result uppercased.
func (s *Slug) Upper(t string) string {
	return strings.ToUpper(s.Make(t))
}

// IsValid reports whether t is already a canonical slug for this
// configuration: non-empty and left unchanged by Make. This means it
// contains only letters, digits and single separators, without leading or
// trailing separators, and (if a limit is set) is within MaxLength.
func (s *Slug) IsValid(t string) bool {
	return t != "" && s.Make(t) == t
}

// MakeUnique returns a slug of t that is unique according to exists: if the
// base slug is already taken it appends the separator and an incrementing
// number ("-2", "-3", ...) until exists reports the candidate as free. When
// MaxLength is set the base is shortened on a word boundary to make room for
// the suffix. A nil exists function disables the uniqueness check.
func (s *Slug) MakeUnique(t string, exists func(string) bool) string {
	base := s.Make(t)
	if exists == nil || !exists(base) {
		return base
	}

	for i := 2; ; i++ {
		if candidate := s.cfg.withSuffix(base, i); !exists(candidate) {
			return candidate
		}
	}
}
