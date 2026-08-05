package slug

import "testing"

// TestPackageParity asserts every package-level helper is exactly its
// default-Slug counterpart, so documentation and behaviour cannot drift
// between the two entry points.
func TestPackageParity(t *testing.T) {
	s := New()
	inputs := []string{
		"Hello World", "co-operate", "email@site.com", "R&D 100%",
		"!!!", "", "MiXeD-CaSe", "北京市 2024",
	}

	for _, in := range inputs {
		if Make(in) != s.Make(in) {
			t.Errorf("Make(%q): package %q != object %q", in, Make(in), s.Make(in))
		}
		if MakeURL(in) != s.MakeURL(in) {
			t.Errorf("MakeURL(%q): package %q != object %q",
				in, MakeURL(in), s.MakeURL(in))
		}
		if Lower(in) != s.Lower(in) {
			t.Errorf("Lower(%q): package %q != object %q", in, Lower(in), s.Lower(in))
		}
		if Upper(in) != s.Upper(in) {
			t.Errorf("Upper(%q): package %q != object %q", in, Upper(in), s.Upper(in))
		}
		if IsValid(in) != s.IsValid(in) {
			t.Errorf("IsValid(%q): package %v != object %v", in, IsValid(in), s.IsValid(in))
		}
	}

	// MakeUnique parity with a simple predicate.
	taken := map[string]bool{"hello-world": true}
	exists := func(x string) bool { return taken[x] }
	if MakeUnique("Hello World", exists) != s.MakeUnique("Hello World", exists) {
		t.Error("MakeUnique parity broken")
	}
}

// TestDefaults documents the zero-option configuration.
func TestDefaults(t *testing.T) {
	s := New()
	if s.cfg.separator != DefaultSeparator {
		t.Errorf("default separator = %q, want %q", s.cfg.separator, DefaultSeparator)
	}
	if s.cfg.maxLen != 0 {
		t.Errorf("default maxLen = %d, want 0", s.cfg.maxLen)
	}
	if s.cfg.fallback != "" {
		t.Errorf("default fallback = %q, want empty", s.cfg.fallback)
	}
}
