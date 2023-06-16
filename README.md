[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/slug)](https://goreportcard.com/report/github.com/goloop/slug) [![License](https://img.shields.io/badge/license-MIT-brightgreen)](https://github.com/goloop/slug/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/slug) [![Stay with Ukraine](https://img.shields.io/static/v1?label=Stay%20with&message=Ukraine%20♥&color=ffd700&labelColor=007acc&style=flat)](https://u24.gov.ua/)


# slug

Package slug generate slug from Unicode string, URL-friendly slugify with
multiple languages support.


## Installation

To install this module use `go get` as:

```
$ go get -u github.com/goloop/slug
```

## Quick Start

To use this module import it as: `github.com/goloop/slug`

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

## Functions

- **Lower**(t string) string

  Lower returns slug in lowercase.

- **Make**(t string) string

  Make returns slug from string.

- **Upper**(t string) string

  Upper returns slug in uppercase.

- **Version**() string

  Version returns the version of the module.

- **New*() *Slug

  New returns pointer to Slug.


## Methods of Slug

- **Lang**(l string) *Slug

  Lang sets the type of language features to use during slugify.

- **Lower**(t string) string

  Lower returns slug in lowercase.

- **Make**(t string) string

  Make returns slug from string.

- **Upper**(t string) string

  Upper returns slug in uppercase.
