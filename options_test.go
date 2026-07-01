package slug

import (
	"testing"

	"github.com/goloop/slug/v2/lang"
)

// TestOptionsCompose checks that several options apply together and that each
// one is reflected in the resolved config.
func TestOptionsCompose(t *testing.T) {
	s := New(
		WithLang(lang.UK),
		WithSeparator("_"),
		WithMaxLength(20),
		WithFallback("post"),
	)

	if s.cfg.lang != lang.UK {
		t.Errorf("lang = %q, want %q", s.cfg.lang, lang.UK)
	}
	if s.cfg.separator != "_" {
		t.Errorf("separator = %q, want _", s.cfg.separator)
	}
	if s.cfg.maxLen != 20 {
		t.Errorf("maxLen = %d, want 20", s.cfg.maxLen)
	}
	if s.cfg.fallback != "post" {
		t.Errorf("fallback = %q, want post", s.cfg.fallback)
	}

	if got, want := s.Make("Привіт, світ!"), "Pryvit_svit"; got != want {
		t.Errorf("Make = %q, want %q", got, want)
	}
}

// TestWithMaxLengthGuard ensures non-positive limits are treated as
// "unlimited" rather than silently truncating everything to nothing.
func TestWithMaxLengthGuard(t *testing.T) {
	for _, n := range []int{0, -1, -100} {
		s := New(WithMaxLength(n))
		if s.cfg.maxLen != 0 {
			t.Errorf("WithMaxLength(%d): maxLen = %d, want 0", n, s.cfg.maxLen)
		}
		if got := s.Make("hello world"); got != "hello-world" {
			t.Errorf("WithMaxLength(%d): Make = %q, want hello-world", n, got)
		}
	}
}

// TestWithLangUnknown documents that an unknown code degrades to neutral
// transliteration instead of failing.
func TestWithLangUnknown(t *testing.T) {
	unknown := New(WithLang("zz-not-a-lang"))
	none := New()
	in := "Привіт Світ"
	if unknown.Make(in) != none.Make(in) {
		t.Errorf("unknown language should behave like None: %q vs %q",
			unknown.Make(in), none.Make(in))
	}
}

// TestOptionsIndependent guards against options sharing mutable state: two
// Slugs built from overlapping options must not affect each other.
func TestOptionsIndependent(t *testing.T) {
	a := New(WithSeparator("_"))
	b := New(WithSeparator("."))
	if a.Make("x y") != "x_y" || b.Make("x y") != "x.y" {
		t.Errorf("options leaked between instances: a=%q b=%q",
			a.Make("x y"), b.Make("x y"))
	}
}
