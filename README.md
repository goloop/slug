[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/slug)](https://goreportcard.com/report/github.com/goloop/slug) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/slug/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/slug) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffD700&labelColor=0057B8&style=flat)](https://u24.gov.ua/)


# slug

Package slug generates URL-friendly slugs from Unicode strings with support for multiple languages. It provides both simple functions and a more configurable object-oriented approach.

## Features

- Clean, URL-safe slug generation.
- Multi-language support through transliteration [t13n](https://github.com/goloop/t13n).
- Customizable character replacement.
- Thread-safe operations.
- High performance.
- No external dependencies for core functionality.

## Installation

```bash
go get -u github.com/goloop/slug
```

## Quick Start

To use this module import it as: `github.com/goloop/slug`

```go
package main

import (
    "fmt"
    "github.com/goloop/slug"
)

func main() {
    // Simple usage
    fmt.Println(slug.Make("Hello 世界"))
    // Output: Hello-Shi-Jie

    // With language-specific settings
    s := slug.New()
    fmt.Println(s.Lang("uk").Make("Привіт Світ"))
    // Output: Pryvit-Svit
}
```

### Conversion functions

#### Fast conversion.

Use the `Make` method to convert a string to slug.

```go
package main

import (
	"fmt"

	"github.com/goloop/slug"
)

func main() {
	// Simple generate slug from the string.
	s := slug.Make("Hello 世界")
	h := "https://example.com/"

	fmt.Printf("%s%s\n", h, s)
	// Output: https://example.com/Hello-Shi-Jie
}
```

## API Reference

### Functions

- `Make(t string) string` - Generate a slug from a string.
- `Lower(t string) string` - Generate a lowercase slug.
- `Upper(t string) string` - Generate an uppercase slug.
- `New() *Slug` - Create a new Slug instance for advanced configuration.

### Slug Methods

- `Lang(l string) *Slug` - Set language for transliteration.
- `Make(t string) string` - Generate a slug.
- `Lower(t string) string` - Generate a lowercase slug.
- `Upper(t string) string` - Generate an uppercase slug.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
