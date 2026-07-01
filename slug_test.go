package slug

import (
	"strings"
	"sync"
	"testing"
	"unicode"

	"github.com/goloop/slug/v2/lang"
)

// TestMake covers the core pipeline with the default configuration: the
// dropped-separator bugs, symbol words, apostrophes and empty input.
func TestMake(t *testing.T) {
	tests := []struct {
		name, in, want string
	}{
		// Hyphens and punctuation must act as separators, not vanish.
		{"literal hyphen kept", "co-operate", "co-operate"},
		{"collapsed hyphens", "Hello---World", "Hello-World"},
		{"punctuation splits", "word!!!word", "word-word"},
		{"slash splits", "a/b", "a-b"},
		{"dots split", "x.y.z", "x-y-z"},
		{"underscore splits", "foo_bar", "foo-bar"},
		{"email keeps words", "email@site.com", "email-at-site-com"},

		// Meaningful symbols become spaced words.
		{"ampersand", "r&d", "r-and-d"},
		{"percent", "100%", "100-pct"},
		{"percent glued", "50%off", "50-pct-off"},
		{"hash", "c#", "c-sharp"},

		// Apostrophes join the surrounding letters.
		{"apostrophe its", "it's", "its"},
		{"apostrophe dont", "don't", "dont"},
		{"apostrophe name", "O'Brien", "OBrien"},
		{"apostrophe chain", "rock'n'roll", "rocknroll"},

		// Case is preserved by Make.
		{"case preserved", "Hello World", "Hello-World"},
		{"unicode folded", "Hellö Wörld", "Hello-World"},

		// Trimming and empty results.
		{"trim edges", "\tHellö \t Wörld\n ", "Hello-World"},
		{"empty", "", ""},
		{"only punct", "!!!", ""},
		{"only spaces", "   ", ""},
		{"only hyphens", "---", ""},
		{"plus signs", "C++", "C"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Make(tt.in); got != tt.want {
				t.Errorf("Make(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

// TestLang checks regional transliteration, including the Ukrainian
// apostrophe, hieroglyphs and mixed punctuation.
func TestLang(t *testing.T) {
	uk := New(WithLang(lang.UK))
	tests := []struct {
		in, want string
	}{
		{"Starlink Ілона Маска відкриє офіс в Україні",
			"Starlink-Ilona-Maska-vidkryie-ofis-v-Ukraini"},
		{"Привіт, світ!", "Pryvit-svit"},
		{"п'ять", "piat"},
		{"м'ясо", "miaso"},
		{"об'єкт", "obiekt"},
		{"你好世界", "Ni-Hao-Shi-Jie"},
		{"[^你好世界$]", "Ni-Hao-Shi-Jie"},
		{"This & that", "This-and-that"},
	}

	for _, tt := range tests {
		if got := uk.Make(tt.in); got != tt.want {
			t.Errorf("uk.Make(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

// TestCase verifies Lower/Upper are exactly the cased forms of Make and that
// the package-level helpers agree with a default Slug.
func TestCase(t *testing.T) {
	inputs := []string{"Hello World", "Ünïcödé café", "Привіт Світ", "MiXeD-CaSe"}
	s := New()
	for _, in := range inputs {
		if got, want := Lower(in), strings.ToLower(s.Make(in)); got != want {
			t.Errorf("Lower(%q) = %q, want %q", in, got, want)
		}
		if got, want := Upper(in), strings.ToUpper(s.Make(in)); got != want {
			t.Errorf("Upper(%q) = %q, want %q", in, got, want)
		}
		if strings.IndexFunc(Lower(in), unicode.IsUpper) >= 0 {
			t.Errorf("Lower(%q) = %q contains uppercase", in, Lower(in))
		}
		// Package helpers must match the default Slug exactly.
		if Make(in) != s.Make(in) {
			t.Errorf("package Make disagrees with default Slug on %q", in)
		}
	}
}

// TestFallback confirms the empty-result fallback fires only when needed.
func TestFallback(t *testing.T) {
	s := New(WithFallback("post"))
	if got := s.Make("!!!"); got != "post" {
		t.Errorf("empty input: got %q, want %q", got, "post")
	}
	if got := s.Make(""); got != "post" {
		t.Errorf("blank input: got %q, want %q", got, "post")
	}
	if got := s.Make("Hello"); got != "Hello" {
		t.Errorf("non-empty input must ignore fallback: got %q", got)
	}
}

// TestSeparator checks a custom separator, including the empty one, and that
// the separator never leaks into a leading/trailing position.
func TestSeparator(t *testing.T) {
	tests := []struct {
		sep, in, want string
	}{
		{"_", "Hello World", "Hello_World"},
		{"_", "  a  b  ", "a_b"},
		{"", "Hello World", "HelloWorld"},
		{".", "a/b/c", "a.b.c"},
		{"~", "one two three", "one~two~three"},
	}
	for _, tt := range tests {
		s := New(WithSeparator(tt.sep))
		if got := s.Make(tt.in); got != tt.want {
			t.Errorf("sep %q, Make(%q) = %q, want %q", tt.sep, tt.in, got, tt.want)
		}
	}
}

// TestMaxLength verifies word-boundary truncation and the hard cut for a
// single oversized word, and that the limit is never exceeded.
func TestMaxLength(t *testing.T) {
	tests := []struct {
		n        int
		in, want string
	}{
		{11, "hello world foo bar", "hello-world"},
		{11, "aa bb cc dd", "aa-bb-cc-dd"},
		{5, "hello world", "hello"},
		{11, "supercalifragilistic", "supercalifr"}, // single word, hard cut
		{3, "ab cd", "ab"},
		{1, "hello", "h"},
	}
	for _, tt := range tests {
		s := New(WithMaxLength(tt.n))
		got := s.Make(tt.in)
		if got != tt.want {
			t.Errorf("max %d, Make(%q) = %q, want %q", tt.n, tt.in, got, tt.want)
		}
		if len(got) > tt.n {
			t.Errorf("result %q exceeds maxLen %d", got, tt.n)
		}
	}
}

// TestIsValid checks the fixed-point definition of a canonical slug.
func TestIsValid(t *testing.T) {
	valid := []string{"hello-world", "Hello-World", "a", "a1-b2", "abc123"}
	invalid := []string{"", "hello--world", "-hello", "hello-", "hello world",
		"hello_world", "café", "hello!"}

	for _, s := range valid {
		if !IsValid(s) {
			t.Errorf("IsValid(%q) = false, want true", s)
		}
	}
	for _, s := range invalid {
		if IsValid(s) {
			t.Errorf("IsValid(%q) = true, want false", s)
		}
	}

	// A custom separator changes what counts as valid.
	u := New(WithSeparator("_"))
	if !u.IsValid("hello_world") || u.IsValid("hello-world") {
		t.Error("IsValid must respect the configured separator")
	}
}

// TestMakeUnique checks suffix allocation, the nil predicate and the
// interaction with MaxLength.
func TestMakeUnique(t *testing.T) {
	taken := map[string]bool{"post": true, "post-2": true, "post-3": false}
	got := MakeUnique("post", func(s string) bool { return taken[s] })
	if got != "post-3" {
		t.Errorf("got %q, want post-3", got)
	}

	// nil predicate disables the check.
	if got := MakeUnique("post", nil); got != "post" {
		t.Errorf("nil exists: got %q, want post", got)
	}

	// First candidate free.
	if got := MakeUnique("fresh", func(string) bool { return false }); got != "fresh" {
		t.Errorf("got %q, want fresh", got)
	}

	// MaxLength must still hold for the suffixed slug.
	s := New(WithMaxLength(6))
	uniq := s.MakeUnique("hello world", func(x string) bool { return x == "hello" })
	if len(uniq) > 6 {
		t.Errorf("unique slug %q exceeds maxLen 6", uniq)
	}
	if uniq == "hello" {
		t.Error("MakeUnique returned the taken base")
	}
}

// TestInvariants is a property test: for a broad set of inputs the output of
// Make must always be a canonical slug — allowed characters only, no leading,
// trailing or doubled separators — and Make must be idempotent (a fixed
// point). Idempotence is the single strongest reliability guarantee here.
func TestInvariants(t *testing.T) {
	inputs := []string{
		"", " ", "-", "--", "a", "A", "1",
		"co-operate", "word!!!word", "  多  spaces  ",
		"email@site.com", "r&d", "100%", "c#", "a_b_c",
		"it's a test", "О'Генрі", "北京市 2024",
		"Ünïcödé café ☕ déjà", "!!!@@@###", "MiXeD.CaSe/Path",
		"tab\tnew\nline", "trailing---", "---leading",
		strings.Repeat("a-", 50), strings.Repeat("世 ", 40),
		string([]byte{0xff, 0xfe, 0x00, 'a', 0xff, 'b'}),
	}

	seps := []string{"-", "_", "", "~"}
	for _, sep := range seps {
		s := New(WithSeparator(sep))
		for _, in := range inputs {
			out := s.Make(in)
			assertCanonical(t, out, sep, in)

			// Idempotence: re-slugging a slug changes nothing.
			if again := s.Make(out); again != out {
				t.Errorf("not idempotent (sep %q): Make(%q)=%q, Make again=%q",
					sep, in, out, again)
			}

			// A non-empty slug must validate under the same config.
			if out != "" && !s.IsValid(out) {
				t.Errorf("IsValid(%q)=false for output of %q (sep %q)", out, in, sep)
			}
		}
	}
}

// assertCanonical fails if s is not a well-formed slug for the given
// separator: only [A-Za-z0-9] and the separator, and no leading, trailing or
// doubled separators when the separator is non-empty.
func assertCanonical(t *testing.T, s, sep, from string) {
	t.Helper()
	if s == "" {
		return
	}

	body := s
	if sep != "" {
		if strings.HasPrefix(s, sep) || strings.HasSuffix(s, sep) {
			t.Errorf("output %q (from %q) has a boundary separator", s, from)
		}
		if strings.Contains(s, sep+sep) {
			t.Errorf("output %q (from %q) has a doubled separator", s, from)
		}
		body = strings.ReplaceAll(s, sep, "")
	}

	for _, r := range body {
		if !(r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			t.Errorf("output %q (from %q) has illegal rune %q", s, from, r)
			return
		}
	}
}

// TestChunkStitching stresses the t13n parallelization path (inputs longer
// than its single-threaded threshold): a long hieroglyph run must produce
// exactly one token per glyph with no words glued or split at chunk borders.
func TestChunkStitching(t *testing.T) {
	const n = 400 // well above the ~256-rune threshold
	in := strings.Repeat("世 ", n)
	out := Make(in)

	want := strings.Repeat("Shi-", n-1) + "Shi"
	if out != want {
		t.Errorf("chunk stitching broke: got %q...", out[:min(60, len(out))])
	}
	if c := strings.Count(out, "Shi"); c != n {
		t.Errorf("expected %d tokens, got %d", n, c)
	}
	assertCanonical(t, out, "-", in)
}

// TestDeterminism ensures repeated calls (including the parallel path) yield
// identical results — no data races leaking through the shared t13n engine.
func TestDeterminism(t *testing.T) {
	in := strings.Repeat("Привіт, світ! ", 60) // > 256 runes
	s := New(WithLang(lang.UK))
	first := s.Make(in)
	for i := 0; i < 20; i++ {
		if got := s.Make(in); got != first {
			t.Fatalf("non-deterministic output on iteration %d", i)
		}
	}
	assertCanonical(t, first, "-", in)
}

// TestConcurrent shares one *Slug across many goroutines; run with -race it
// proves a Slug is safe for concurrent use.
func TestConcurrent(t *testing.T) {
	s := New(WithLang(lang.UK), WithMaxLength(40), WithFallback("x"))
	inputs := []string{"Привіт світ", "Hello World", "!!!", "co-operate",
		strings.Repeat("北京 ", 200), "r&d 100%"}

	want := make([]string, len(inputs))
	for i, in := range inputs {
		want[i] = s.Make(in)
	}

	var wg sync.WaitGroup
	for g := 0; g < 32; g++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, in := range inputs {
				if got := s.Make(in); got != want[i] {
					t.Errorf("concurrent mismatch for %q: got %q, want %q",
						in, got, want[i])
				}
			}
		}()
	}
	wg.Wait()
}

// TestNoPanic throws hostile inputs at Make: invalid UTF-8, NUL bytes, every
// single byte value and a very long random-ish string.
func TestNoPanic(t *testing.T) {
	var every strings.Builder
	for b := 0; b < 256; b++ {
		every.WriteByte(byte(b))
	}

	inputs := []string{
		"", string([]byte{0xff, 0xfe, 0xfd}), "\x00\x00\x00",
		every.String(), strings.Repeat(every.String(), 8),
		strings.Repeat("🙂", 1000), strings.Repeat("a\xffб ", 300),
	}
	for _, in := range inputs {
		func() {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Make panicked on %q: %v", in, r)
				}
			}()
			out := Make(in)
			assertCanonical(t, out, "-", in)
		}()
	}
}

// FuzzMake asserts the canonical-slug invariants and idempotence for
// arbitrary inputs discovered by the fuzzer.
func FuzzMake(f *testing.F) {
	seeds := []string{
		"", "Hello World", "co-operate", "email@site.com", "п'ять",
		"你好世界", "100% sure", string([]byte{0xff, 0x00}), "---",
		strings.Repeat("世 ", 300),
	}
	for _, s := range seeds {
		f.Add(s)
	}

	s := New(WithLang(lang.UK))
	f.Fuzz(func(t *testing.T, in string) {
		out := s.Make(in)
		assertCanonical(t, out, "-", in)

		if again := s.Make(out); again != out {
			t.Errorf("not idempotent: Make(%q)=%q, again=%q", in, out, again)
		}
		if out != "" && !s.IsValid(out) {
			t.Errorf("IsValid(%q)=false for output of %q", out, in)
		}
	})
}
