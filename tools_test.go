package slug

import (
	"testing"

	"github.com/goloop/t13n/lang"
)

// TestSlugRules checks the per-rune transliteration callback in isolation:
// apostrophes are dropped, meaningful symbols become spaced words, and
// everything else is passed through untouched.
func TestSlugRules(t *testing.T) {
	tests := []struct {
		name string
		in   lang.TransState
		want string
	}{
		{"apostrophe dropped", lang.TransState{Curr: '\'', IsApostrophe: true}, ""},
		{"backtick apostrophe", lang.TransState{Curr: '`', IsApostrophe: true}, ""},
		{"lone quote passes", lang.TransState{Curr: '\'', Value: "'"}, "'"},
		{"at", lang.TransState{Curr: '@', Value: "@"}, symAt},
		{"and", lang.TransState{Curr: '&', Value: "&"}, symAnd},
		{"sharp", lang.TransState{Curr: '#', Value: "#"}, symSharp},
		{"pct", lang.TransState{Curr: '%', Value: "%"}, symPct},
		{"letter passthrough", lang.TransState{Curr: 'i', Value: "i"}, "i"},
		{"hieroglyph passthrough", lang.TransState{Curr: '世', Value: "Shi "}, "Shi "},
		{"space passthrough", lang.TransState{Curr: ' ', Value: " "}, " "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, seek, ok := slugRules(tt.in)
			if got != tt.want {
				t.Errorf("value: got %q, want %q", got, tt.want)
			}
			if seek != 0 {
				t.Errorf("seek: got %d, want 0", seek)
			}
			if !ok {
				t.Error("ok: got false, want true")
			}
		})
	}
}

// TestTokenize verifies the whitelist splitter keeps only letter/digit runs
// and never yields empty tokens, regardless of surrounding punctuation.
func TestTokenize(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"", nil},
		{"!!!", nil},
		{"---", nil},
		{"abc", []string{"abc"}},
		{"a-b-c", []string{"a", "b", "c"}},
		{"--a--b--", []string{"a", "b"}},
		{"a1b2", []string{"a1b2"}},
		{" leading", []string{"leading"}},
		{"trailing ", []string{"trailing"}},
		// tokenize sees only the ASCII output of t13n; multibyte runes never
		// reach it. Fed raw UTF-8, it keeps just the ASCII bytes in between.
		{"Ünïcödé", []string{"n", "c", "d"}},
		{"a\x00b", []string{"a", "b"}},
	}

	for _, tt := range tests {
		got := tokenize(tt.in)
		if len(got) != len(tt.want) {
			t.Errorf("tokenize(%q) = %q, want %q", tt.in, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("tokenize(%q)[%d] = %q, want %q", tt.in, i, got[i], tt.want[i])
			}
			if got[i] == "" {
				t.Errorf("tokenize(%q) produced an empty token", tt.in)
			}
		}
	}
}

// TestAssemble exercises separator joining, the fallback and word-boundary
// truncation directly on the config, without the transliteration step.
func TestAssemble(t *testing.T) {
	tests := []struct {
		name   string
		cfg    config
		tokens []string
		want   string
	}{
		{"basic", config{separator: "-"}, []string{"a", "b", "c"}, "a-b-c"},
		{"custom sep", config{separator: "_"}, []string{"a", "b"}, "a_b"},
		{"empty sep", config{separator: ""}, []string{"a", "b"}, "ab"},
		{"empty -> fallback", config{separator: "-", fallback: "n-a"}, nil, "n-a"},
		{"empty -> empty", config{separator: "-"}, nil, ""},
		{"maxlen boundary", config{separator: "-", maxLen: 3}, []string{"aa", "bb", "cc"}, "aa"},
		{"maxlen exact", config{separator: "-", maxLen: 5}, []string{"aa", "bb", "cc"}, "aa-bb"},
		{"maxlen hard cut", config{separator: "-", maxLen: 4}, []string{"abcdef"}, "abcd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.cfg.assemble(tt.tokens); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

// TestWithSuffix checks the uniqueness suffix, including its interaction
// with MaxLength (dropping whole words, then hard-cutting a single word).
func TestWithSuffix(t *testing.T) {
	tests := []struct {
		name string
		cfg  config
		base string
		n    int
		want string
	}{
		{"plain", config{separator: "-"}, "post", 2, "post-2"},
		{"empty base", config{separator: "-"}, "", 2, "2"},
		{"fits", config{separator: "-", maxLen: 10}, "post", 3, "post-3"},
		{"drop word", config{separator: "-", maxLen: 6}, "aa-bb-cc", 2, "aa-2"},
		{"hard cut word", config{separator: "-", maxLen: 6}, "hello", 2, "hell-2"},
		{"empty sep cut", config{separator: "", maxLen: 5}, "hello", 9, "hell9"},
		{"no room", config{separator: "-", maxLen: 2}, "hello", 2, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.cfg.withSuffix(tt.base, tt.n)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
			if tt.cfg.maxLen > 0 && len(got) > tt.cfg.maxLen {
				t.Errorf("result %q exceeds maxLen %d", got, tt.cfg.maxLen)
			}
		})
	}
}

// TestIsAlnum guards the ASCII classifier used by tokenize.
func TestIsAlnum(t *testing.T) {
	for b := 0; b < 256; b++ {
		want := b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
		if got := isAlnum(byte(b)); got != want {
			t.Errorf("isAlnum(%d) = %v, want %v", b, got, want)
		}
	}
	// Bytes that look tempting but must be rejected.
	for _, b := range []byte{'-', '_', ' ', '.', '/', '@', 0x00, 0xFF} {
		if isAlnum(b) {
			t.Errorf("isAlnum(%q) = true, want false", b)
		}
	}
}
