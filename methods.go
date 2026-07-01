package slug

// defaultSlug backs the package-level helpers. It is created once and never
// mutated, so the helpers are safe for concurrent use.
var defaultSlug = New()

// Make returns the slug of t using the default configuration (neutral
// transliteration, "-" separator, original letter case preserved).
// It is shorthand for New().Make(t).
func Make(t string) string {
	return defaultSlug.Make(t)
}

// Lower is Make with the result lowercased.
func Lower(t string) string {
	return defaultSlug.Lower(t)
}

// Upper is Make with the result uppercased.
func Upper(t string) string {
	return defaultSlug.Upper(t)
}

// IsValid reports whether t is already a canonical slug under the default
// configuration (see (*Slug).IsValid).
func IsValid(t string) bool {
	return defaultSlug.IsValid(t)
}

// MakeUnique returns a slug of t made unique via exists (see
// (*Slug).MakeUnique), using the default configuration.
func MakeUnique(t string, exists func(string) bool) string {
	return defaultSlug.MakeUnique(t, exists)
}

// TryMakeUnique returns a unique slug of t and reports success (see
// (*Slug).TryMakeUnique), using the default configuration.
func TryMakeUnique(t string, exists func(string) bool, maxTries int) (string, bool) {
	return defaultSlug.TryMakeUnique(t, exists, maxTries)
}
