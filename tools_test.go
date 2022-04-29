package slug

import (
	"testing"

	"github.com/goloop/t13n/lang"
)

// TestSlugRules tests slugRules function.
func TestSlugRules(t *testing.T) {
	var tests = []struct {
		item  lang.TransState
		value string
		seek  int
		ok    bool
	}{
		{
			item:  lang.TransState{Curr: ' ', Value: " "},
			value: "-",
			seek:  0,
			ok:    true,
		},
		{
			item:  lang.TransState{Curr: '@', Value: "@"},
			value: "at",
			seek:  0,
			ok:    true,
		},
		{
			item:  lang.TransState{Curr: 'i', Value: "i"},
			value: "i",
			seek:  0,
			ok:    true,
		},
		{
			item:  lang.TransState{Curr: '世', Value: "shi"},
			value: "shi",
			seek:  0,
			ok:    true,
		},
		{
			item:  lang.TransState{Curr: '㹏', Next: '㹐', Value: "jin "},
			value: "jin-",
			seek:  0,
			ok:    true,
		},
	}

	for _, test := range tests {
		value, seek, ok := slugRules(test.item)
		if value != test.value {
			t.Errorf("expected %s but %s", test.value, value)
		}

		if seek != test.seek {
			t.Errorf("expected %d but %d", test.seek, seek)
		}

		if ok != test.ok {
			t.Errorf("expected %t but %t", test.ok, ok)
		}
	}
}
