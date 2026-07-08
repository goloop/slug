# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.2.0]

Minor release: fixes to fallback handling and word boundaries. Behavioural
changes are noted below.

### Fixed
- `WithFallback` is now normalized through the slug pipeline at construction
  time, so the fallback always obeys the canonical-slug and `MaxLength`
  invariants instead of being returned verbatim. A fallback that normalizes to
  nothing is treated as no fallback.
- A visible character with no transliteration (an emoji, an out-of-table
  symbol) now acts as a word boundary instead of vanishing and gluing its
  neighbours, so words separated by such a character are no longer merged.
  Invisible format characters (ZWJ/ZWNJ/soft hyphen) still join the
  surrounding letters.
- `New` tolerates `nil` options instead of panicking, which supports the
  common pattern of conditionally assembling an option slice.
- `TryMakeUnique` reports `("", false)` when the input produces an empty base
  slug and no fallback is set, instead of the misleading `("", true)`.

### Changed
- Bumped the transliteration dependency to t13n/v2 v2.1.1 (documentation-only
  release; no change to slug output).

## [2.1.0]

Minor release: fixes to `MakeUnique` and a faster `Lower`/`Upper`. Fully
backward compatible.

### Added
- `TryMakeUnique(t, exists, maxTries) (string, bool)` — a bounded, strict
  variant of `MakeUnique` that reports whether a unique slug was produced.

### Fixed
- `MakeUnique` no longer loops forever when the `exists` predicate always
  reports "taken"; it now gives up after a bounded number of attempts and
  returns the base slug (best-effort). Use `TryMakeUnique` to detect failure.
- `MakeUnique` no longer exceeds `MaxLength`: when even a bare numeric suffix
  would not fit, the search reports failure instead of emitting an
  over-length slug.

### Changed
- Bumped the transliteration dependency to t13n/v2 v2.1.0 (API-compatible;
  no change to slug output).
- `Lower` and `Upper` build the slug directly in the target case in a single
  pass (one fewer allocation).
- `WithSeparator` documents that the separator is not validated and should be
  a URL-safe, non-alphanumeric string.
- `IsValid` documents its interaction with `WithFallback` (it checks that the
  input is itself canonical, not that it would produce a valid slug).
- Removed the stale `benchmarks.txt` snapshot (see `benchmarks_test.go`).

## [2.0.0]

Major release. The module path is now `github.com/goloop/slug/v2` and the
minimum Go version is 1.24. The public API was reworked around functional
options; import paths and some outputs change, so this is a breaking release.

### Added
- Functional options: `WithLang`, `WithSeparator`, `WithMaxLength`,
  `WithFallback`, applied through `New(...)`.
- `MaxLength` support with truncation on a word boundary (a word is never
  split unless it alone exceeds the limit).
- `Fallback` value returned when the input yields an empty slug.
- Custom `Separator` (defaults to `-`; may be empty).
- `IsValid` reports whether a string is already a canonical slug.
- `MakeUnique` appends an incrementing suffix until a predicate reports the
  slug as free, respecting `MaxLength`.
- Runnable examples, property/invariant tests, a fuzz target and coverage
  for the `lang` subpackage.

### Changed
- A hyphen in the input is now kept as a separator instead of being deleted
  (`co-operate` -> `co-operate`, previously `cooperate`).
- Punctuation between words now separates them instead of gluing them
  (`word!!!word` -> `word-word`, `email@site.com` -> `email-at-site-com`).
- Meaningful symbols become standalone words: `r&d` -> `r-and-d`,
  `100%` -> `100-pct`, `c#` -> `c-sharp`.
- `Make` preserves the input letter case; use `Lower`/`Upper` to force one.
- A `Slug` is immutable after `New` and safe for concurrent use.

### Removed
- The fluent `(*Slug).Lang(string) *Slug` setter (use `WithLang`).

### Fixed
- Lost input hyphens and dropped word separators that silently merged words.
- Per-rune allocation of the ignore table and the two regular-expression
  passes were replaced by a single linear pass.
