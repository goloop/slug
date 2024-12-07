// Package slug provides functionality to generate URL-friendly slugs
// from Unicode strings with support for multiple languages.
//
// This package is designed to create clean, URL-safe strings from
// any input text, with special handling for various languages and
// character sets. It supports:
//
//   - Basic Latin characters (a-z, A-Z).
//   - Numbers (0-9).
//   - Special character conversion (spaces, underscores to hyphens).
//   - Multi-language support through transliteration.
//   - Custom character replacement rules.
//   - Case transformation (upper/lower).
//
// Usage:
//
// Basic usage with the package-level functions:
//
//	slug.Make("Hello World")  // returns "Hello-World"
//	slug.Lower("Hello World") // returns "hello-world"
//	slug.Upper("Hello World") // returns "HELLO-WORLD"
//
// Using the Slug type for more control:
//
//	s := slug.New()
//	s.Lang("uk").Make("Привіт Світ") // returns transliterated version
//
// The package is thread-safe and can be used in concurrent applications.
package slug
