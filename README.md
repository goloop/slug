[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/slug)](https://goreportcard.com/report/github.com/goloop/slug) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/slug/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/slug/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# slug

Package slug generates URL-friendly slugs from Unicode strings in any
language. It offers simple package-level functions and a configurable,
immutable `Slug` type built from functional options.

## Features

- Clean, URL-safe slug generation from any Unicode text.
- Multi-language transliteration via [t13n](https://github.com/goloop/t13n).
- Punctuation and symbols handled predictably (never silently merges words).
- Configurable separator, maximum length and empty-input fallback.
- Uniqueness helper and slug validation.
- Immutable and safe for concurrent use.

## Installation

```bash
go get github.com/goloop/slug/v2
```

## Quick Start

```go
package main

import (
	"fmt"

	"github.com/goloop/slug/v2"
)

func main() {
	fmt.Println(slug.Make("Hello World"))  // Hello-World
	fmt.Println(slug.Lower("Hello World")) // hello-world
	fmt.Println(slug.Upper("Hello World")) // HELLO-WORLD

	// Punctuation separates words; it is never dropped silently.
	fmt.Println(slug.Make("co-operate"))     // co-operate
	fmt.Println(slug.Make("email@site.com")) // email-at-site-com
	fmt.Println(slug.Make("R&D 100%"))       // R-and-D-100-pct
}
```

## Configuration

Build a `Slug` with `New` and functional options for full control:

```go
s := slug.New(
	slug.WithLang(lang.UK),    // regional transliteration
	slug.WithSeparator("_"),   // custom word separator (default "-")
	slug.WithMaxLength(60),    // cut on a word boundary
	slug.WithFallback("post"), // value for an empty result
)

s.Make("Привіт, світ!") // Pryvit_svit
```

A `Slug` is immutable after `New`, so one value may be shared safely across
goroutines.

### Behaviour notes

- `Make` preserves the input letter case; use `Lower` or `Upper` to force one.
- Meaningful symbols become words: `@` → `at`, `&` → `and`, `#` → `sharp`,
  `%` → `pct`.
- A true apostrophe between letters joins them: `it's` → `its`.
- Any other non-alphanumeric character acts as a word separator; leading,
  trailing and repeated separators are removed.

### Uniqueness and validation

```go
// Ensure a unique slug against your storage.
s.MakeUnique("hello world", func(candidate string) bool {
	return exists(candidate) // your lookup; returns true if taken
}) // hello-world, hello-world-2, ...

// Strict variant: bounded, and reports whether a unique slug was found.
if slug, ok := s.TryMakeUnique("hello world", exists, 100); ok {
	use(slug)
}

slug.IsValid("hello-world") // true
slug.IsValid("hello world") // false
```

## API Reference

### Package functions

- `Make(t string) string` — slug with the default configuration.
- `Lower(t string) string` — lowercase slug.
- `Upper(t string) string` — uppercase slug.
- `IsValid(t string) bool` — whether `t` is already a canonical slug.
- `MakeUnique(t string, exists func(string) bool) string` — unique slug.
- `TryMakeUnique(t string, exists func(string) bool, maxTries int) (string, bool)` — bounded unique slug.

### Options

- `WithLang(code string) Option` — transliteration language (see `lang`).
- `WithSeparator(sep string) Option` — word separator (default `-`).
- `WithMaxLength(n int) Option` — length limit, cut on a word boundary.
- `WithFallback(s string) Option` — value for an empty result.

### Slug

- `New(opts ...Option) *Slug` — create a configured, immutable slug maker.
- `(*Slug).Make / Lower / Upper / IsValid / MakeUnique` — as above.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
