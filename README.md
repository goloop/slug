[//]: # (!!!Don't modify the README.md, use `make readme` to generate it!!!)


[![Go Report Card](https://goreportcard.com/badge/github.com/goloop/slug)](https://goreportcard.com/report/github.com/goloop/slug) [![License](https://img.shields.io/badge/license-BSD-blue)](https://github.com/goloop/slug/blob/master/LICENSE) [![License](https://img.shields.io/badge/godoc-YES-green)](https://godoc.org/github.com/goloop/slug)

*Version: v0.0.1-alpha*

# slug

Package slug generate slug from Unicode string, URL-friendly slugify with
multiple languages support.


## Installation

To install this module use `go get` as:

    $ go get -u github.com/goloop/slug

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

	// Output:
	//      hello-shi-jie
}
```

## Usage

#### func  Lower

    func Lower(t string) string

Lower returns slug in lowercase.

#### func  Make

    func Make(t string) string

Make returns slug from string.

#### func  Version

    func Version() string

Version returns the version of the module.

#### type Slug

    type Slug struct {
    }


Slug is the slug constructor.

#### func  New

    func New() *Slug

New retursn pointer to Slug.

#### func (*Slug) Lang

    func (s *Slug) Lang(l string) *Slug

Lang sets the type of language features to use during slugify.

#### func (*Slug) Lower

    func (s *Slug) Lower(t string) string

Lower returns slug in lowercase.

#### func (*Slug) Make

    func (s *Slug) Make(t string) string

Make returns slug from string.
