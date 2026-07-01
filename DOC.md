# slug — reference

The full reference for the `slug` package: the mental model, the package
functions, the configurable `Slug` type, the options, the normalization rules,
uniqueness, validation and practical recipes.

Ukrainian version: **[DOC.UK.md](DOC.UK.md)**.

## Contents

- [Mental model](#mental-model)
- [Package functions](#package-functions)
- [The Slug type and options](#the-slug-type-and-options)
- [Normalization rules](#normalization-rules)
- [Transliteration and language](#transliteration-and-language)
- [Case](#case)
- [Uniqueness](#uniqueness)
- [Validation](#validation)
- [Recipes and tips](#recipes-and-tips)

## Mental model

`slug` turns arbitrary Unicode text into a URL-friendly slug, in any language.
There are two ways to use it:

1. **Package functions** — `slug.Make`, `slug.Lower`, `slug.Upper`, … use the
   default configuration.
2. **A configured `Slug`** — build one with `New` and functional options for a
   custom separator, length limit, transliteration language and empty-input
   fallback.

The guiding principle is that **word boundaries are never silently lost**. A
punctuation mark separates words rather than gluing them together, and a handful
of meaningful symbols become their own words (`@` → `at`, `&` → `and`). So the
result stays readable and reversible in meaning, not a run-together blob.

A `Slug` is **immutable after `New`**, so one value can be shared safely across
goroutines.

```go
import "github.com/goloop/slug/v2"
```

## Package functions

```go
func Make(t string) string
func Lower(t string) string
func Upper(t string) string
func IsValid(t string) bool
func MakeUnique(t string, exists func(string) bool) string
func TryMakeUnique(t string, exists func(string) bool, maxTries int) (string, bool)
```

Each uses the default configuration (separator `-`, no length limit, no
transliteration language):

```go
slug.Make("Hello World")  // "Hello-World"
slug.Lower("Hello World") // "hello-world"
slug.Upper("Hello World") // "HELLO-WORLD"
```

## The Slug type and options

```go
func New(opts ...Option) *Slug

func WithLang(code string) Option
func WithSeparator(sep string) Option // default "-" (DefaultSeparator)
func WithMaxLength(n int) Option
func WithFallback(s string) Option
```

`New` returns a configured, immutable maker exposing the same methods as the
package functions (`Make`, `Lower`, `Upper`, `IsValid`, `MakeUnique`,
`TryMakeUnique`):

```go
s := slug.New(
    slug.WithLang(lang.UK),    // regional transliteration
    slug.WithSeparator("_"),   // custom word separator (default "-")
    slug.WithMaxLength(60),    // cut on a word boundary
    slug.WithFallback("post"), // value for an empty result
)

s.Make("Привіт, світ!") // "Pryvit_svit"
```

| Option | Effect |
|--------|--------|
| `WithLang(code)`      | choose the transliteration language (see [below](#transliteration-and-language)) |
| `WithSeparator(sep)`  | the word separator (default `-`) |
| `WithMaxLength(n)`    | limit the length, cutting on a word boundary so no word is split |
| `WithFallback(s)`     | the value returned when the slug would otherwise be empty |

## Normalization rules

The normalization is a whitelist: only letters and digits survive as content,
everything else either becomes a word or a separator. The rules are:

- **Meaningful symbols become words:** `@` → `at`, `&` → `and`, `#` → `sharp`,
  `%` → `pct`.
- **A true apostrophe between letters joins them:** `it's` → `its`.
- **Any other non-alphanumeric character is a word separator.** Leading,
  trailing and repeated separators are collapsed and removed.

```go
slug.Make("co-operate")     // "co-operate"
slug.Make("email@site.com") // "email-at-site-com"
slug.Make("R&D 100%")       // "R-and-D-100-pct"
```

This whitelist approach is what guarantees separators are never silently merged:
punctuation always breaks words instead of disappearing.

## Transliteration and language

Non-ASCII letters are transliterated to ASCII. `WithLang` selects the regional
rules; the `slug/v2/lang` subpackage exposes the language codes (for example
`lang.UK`, `lang.EN`):

```go
import "github.com/goloop/slug/v2/lang"

s := slug.New(slug.WithLang(lang.UK))
s.Make("Привіт, світ!") // "Pryvit-svit"
```

Without a language, a general transliteration is applied. Choosing the right
language matters when the same character transliterates differently across
locales.

## Case

`Make` **preserves the input letter case**; use `Lower` or `Upper` to force one:

```go
slug.Make("Hello World")  // "Hello-World"
slug.Lower("Hello World") // "hello-world"
slug.Upper("Hello World") // "HELLO-WORLD"
```

## Uniqueness

```go
func (s *Slug) MakeUnique(t string, exists func(string) bool) string
func (s *Slug) TryMakeUnique(t string, exists func(string) bool, maxTries int) (string, bool)
```

`MakeUnique` appends a numeric suffix until your `exists` predicate reports the
candidate is free:

```go
s.MakeUnique("hello world", func(candidate string) bool {
    return exists(candidate) // your lookup; true if taken
}) // "hello-world", "hello-world-2", ...
```

`TryMakeUnique` is the bounded, strict variant — it tries at most `maxTries`
times and reports whether a unique slug was found, so a hostile or exhausted
namespace can't spin forever:

```go
if slug, ok := s.TryMakeUnique("hello world", exists, 100); ok {
    use(slug)
}
```

## Validation

```go
func IsValid(t string) bool
func (s *Slug) IsValid(t string) bool
```

`IsValid` reports whether a string is already a canonical slug for the
configuration — i.e. `Make` would return it unchanged:

```go
slug.IsValid("hello-world") // true
slug.IsValid("hello world") // false
```

## Recipes and tips

**Build one maker, reuse it.** A `Slug` is immutable and concurrency-safe —
construct it once at startup with your options and share it, rather than passing
options at every call site.

**Bound your uniqueness loop.** Prefer `TryMakeUnique` over `MakeUnique` when the
`exists` check hits a database or a large namespace, so a saturated space returns
`ok=false` instead of looping.

**Give empty input a home.** Set `WithFallback` so text that reduces to nothing
(pure punctuation, emoji-only titles) yields a stable placeholder instead of an
empty string.

**Pick the language deliberately.** Use `WithLang` with the right `lang` code
when your content is in a specific language; the default general transliteration
is a safe fallback, not a locale-aware one.

**Cap the length on a boundary.** `WithMaxLength` cuts at a word boundary, so
slugs stay readable and never end mid-word.
