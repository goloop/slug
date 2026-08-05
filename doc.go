// Package slug generates URL-friendly slugs from Unicode strings in any
// language.
//
// A slug is built in three steps: the input is transliterated to ASCII,
// split into words on every character that is not a letter or a digit, and
// the words are joined with a separator ("-" by default). Leading, trailing
// and repeated separators are removed, so the result contains only letters,
// digits and single separators.
//
// # Package-level helpers
//
// The package functions use a default configuration (neutral
// transliteration, "-" separator, original letter case preserved):
//
//	slug.Make("Hello World")    // "Hello-World"
//	slug.MakeURL("Hello World") // "hello-world"
//	slug.Lower("Hello World")   // "hello-world"
//	slug.Upper("Hello World")   // "HELLO-WORLD"
//
// # Slugs that go into a URL
//
// Make preserves the case transliteration produced, which for a non-Latin
// title means a capital letter in the path:
//
//	slug.Make("Осінній настрій")    // "Osinnii-nastrii"
//	slug.MakeURL("Осінній настрій") // "osinnii-nastrii"
//
// URL paths compare case-sensitively, so one title must not be able to reach
// a resource under two spellings. Use MakeURL where the slug becomes part of
// a URL, or WithLowercase to make lower case the default for every method of
// a configured Slug - MakeUnique and TryMakeUnique build on Make, so only the
// option covers them.
//
// # Configurable slugs
//
// New builds a Slug from functional options for full control:
//
//	s := slug.New(
//		slug.WithLang(lang.UK),   // regional transliteration
//		slug.WithSeparator("_"),  // custom separator
//		slug.WithMaxLength(60),   // cut on a word boundary
//		slug.WithFallback("post"),// value for an empty result
//	)
//	s.Make("Привіт, світ!") // "Pryvit_svit"
//
// A few characters carry meaning and become words: "@" -> "at",
// "&" -> "and", "#" -> "sharp", "%" -> "pct" (e.g. "r&d" -> "r-and-d").
// A true apostrophe between letters is dropped so the letters join
// ("it's" -> "its").
//
// # Concurrency
//
// A Slug is immutable after New, so a single value may be shared and used
// by multiple goroutines concurrently; the package-level helpers are safe
// for concurrent use as well.
package slug
