package slug

import (
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
	return s.make(t, casePreserve)
}

// Lower is Make with the result lowercased.
func (s *Slug) Lower(t string) string {
	return s.make(t, caseLower)
}

// Upper is Make with the result uppercased.
func (s *Slug) Upper(t string) string {
	return s.make(t, caseUpper)
}

// make runs the pipeline once, producing the slug directly in the requested
// case (so Lower/Upper need no extra pass).
func (s *Slug) make(t string, mode caseMode) string {
	raw := t13n.Render(s.cfg.lang, t, slugRules)
	return s.cfg.render(raw, mode)
}

// IsValid reports whether t is already a canonical slug for this
// configuration: non-empty and left unchanged by Make. This means it
// contains only letters, digits and single separators, without leading or
// trailing separators, and (if a limit is set) is within MaxLength.
//
// IsValid checks that t is itself canonical, not that t would produce a valid
// slug. In particular, with WithFallback an input that maps to the fallback
// is not canonical: IsValid("!!!") is false even though Make("!!!") returns
// the (valid) fallback.
func (s *Slug) IsValid(t string) bool {
	return t != "" && s.Make(t) == t
}

// defaultUniqueTries bounds MakeUnique so a misbehaving exists predicate
// cannot loop forever.
const defaultUniqueTries = 1 << 16

// MakeUnique returns a slug of t that is unique according to exists: if the
// base slug is already taken it appends the separator and an incrementing
// number ("-2", "-3", ...) until exists reports the candidate as free. When
// MaxLength is set the base is shortened on a word boundary to keep every
// candidate within the limit. A nil exists function disables the check.
//
// MakeUnique is best-effort: it gives up after a bounded number of attempts
// (so a predicate that always reports "taken" cannot hang) and returns the
// base slug, which may then collide. Use TryMakeUnique when you need to detect
// that no unique slug could be produced.
func (s *Slug) MakeUnique(t string, exists func(string) bool) string {
	if slug, ok := s.TryMakeUnique(t, exists, defaultUniqueTries); ok {
		return slug
	}

	return s.Make(t)
}

// TryMakeUnique returns a unique slug of t and reports whether it succeeded.
// It tests the base slug, then base+separator+2, base+separator+3, ... with
// exists (which must report whether a slug is already taken), trying at most
// maxTries candidates. Every candidate satisfies MaxLength; if the limit is so
// small that no numbered candidate fits, the search stops early. A
// non-positive maxTries uses a default bound; a nil exists succeeds with the
// base slug.
//
// It returns (slug, true) on the first free candidate, or ("", false) if none
// is free within the attempts or no length-valid candidate can be formed.
func (s *Slug) TryMakeUnique(t string, exists func(string) bool, maxTries int) (string, bool) {
	if maxTries <= 0 {
		maxTries = defaultUniqueTries
	}

	base := s.Make(t)
	if exists == nil || !exists(base) {
		return base, true
	}

	// The base counts as the first attempt.
	for n, tries := 2, 1; tries < maxTries; n, tries = n+1, tries+1 {
		candidate, ok := s.cfg.withSuffix(base, n)
		if !ok {
			break
		}
		if !exists(candidate) {
			return candidate, true
		}
	}

	return "", false
}
