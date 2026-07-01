package lang

import (
	"testing"

	t13n "github.com/goloop/t13n/lang"
)

// TestReexport guards that the language codes re-exported here still point at
// the correct upstream values and were not accidentally cross-wired
// (e.g. UK = US). A mismatch would silently apply the wrong regional rules.
func TestReexport(t *testing.T) {
	cases := []struct {
		name      string
		got, want string
	}{
		{"None", None, t13n.None},
		{"UK", UK, t13n.UK},
		{"US", US, t13n.US},
		{"GB", GB, t13n.GB},
		{"DE", DE, t13n.DE},
		{"RU", RU, t13n.RU},
		{"BG", BG, t13n.BG},
		{"FR", FR, t13n.FR},
		{"EN", EN, t13n.EN},
	}

	for _, c := range cases {
		if c.got != c.want {
			t.Errorf("%s = %q, want %q", c.name, c.got, c.want)
		}
	}

	if None != "" {
		t.Errorf("None = %q, want empty string", None)
	}

	// Codes that map to different regional rules must be distinct.
	distinct := []string{UK, US, DE, RU, BG, FR}
	seen := make(map[string]bool, len(distinct))
	for _, code := range distinct {
		if code == "" {
			t.Error("a real language code is empty")
		}
		if seen[code] {
			t.Errorf("duplicate language code %q", code)
		}
		seen[code] = true
	}
}
