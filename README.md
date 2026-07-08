[![deps.dev](https://img.shields.io/badge/deps.dev-insights-4c8dbc)](https://deps.dev/go/github.com%2Fgoloop%2Fslug%2Fv2) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/slug/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://pkg.go.dev/github.com/goloop/slug/v2) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)

# slug

`slug` generates URL-friendly slugs from Unicode strings in any language. It
offers simple package-level functions and a configurable, immutable `Slug` type
built from functional options.

Its guiding rule is that **word boundaries are never silently lost**: a
punctuation mark separates words instead of gluing them together, and a handful
of meaningful symbols become their own words (`@` → `at`, `&` → `and`). The
result stays readable, not a run-together blob.

## Features

- Clean, URL-safe slug generation from any Unicode text.
- Multi-language transliteration.
- Punctuation and symbols handled predictably (never silently merges words).
- Configurable separator, maximum length and empty-input fallback.
- Uniqueness helper and slug validation.
- Immutable and safe for concurrent use.

## Installation

```bash
go get github.com/goloop/slug/v2
```

```go
import "github.com/goloop/slug/v2"
```

Requires Go 1.24 or newer.

## Quick start

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

Build a configured, immutable maker for full control:

```go
import "github.com/goloop/slug/v2/lang"

s := slug.New(
    slug.WithLang(lang.UK),    // regional transliteration
    slug.WithSeparator("_"),   // custom word separator (default "-")
    slug.WithMaxLength(60),    // cut on a word boundary
    slug.WithFallback("post"), // value for an empty result
)

s.Make("Привіт, світ!") // Pryvit_svit
```

## Documentation

- Full reference and recipes: [DOC.md](DOC.md) · [DOC.UK.md](DOC.UK.md)
- Package API: [pkg.go.dev/github.com/goloop/slug/v2](https://pkg.go.dev/github.com/goloop/slug/v2)
- Changes between versions: [CHANGELOG.md](CHANGELOG.md)

## Contributing

Contributions are welcome. Please run `go test ./...`, `go vet ./...` and
`gofmt -l .` before submitting a pull request.

## License

`slug` is released under the MIT License. See [LICENSE](LICENSE).
