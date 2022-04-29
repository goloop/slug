package slug

import (
	"testing"

	"github.com/goloop/slug/lang"
)

// TestObjMake tests Make method.
func TestObjMake(t *testing.T) {
	var tests = []struct {
		value    string
		expected string
	}{
		{
			"Starlink Ілона Маска відкриє офіс в Україні",
			"Starlink-Ilona-Maska-vidkryie-ofis-v-Ukraini",
		},
		{"Hellö Wörld", "Hello-World"},
		{"你好世界", "Ni-Hao-Shi-Jie"},
		{"[^你好世界$]", "Ni-Hao-Shi-Jie"},
		{"This & that", "This-and-that"},
		{"\tHellö \t Wörld\n ", "Hello-World"},
	}

	s := New()
	s.Lang(lang.UK)
	for _, test := range tests {
		if v := s.Make(test.value); v != test.expected {
			t.Errorf("expected %s but %s", test.expected, v)
		}
	}
}

// TestObjLower tests Lower method.
func TestObjLower(t *testing.T) {
	var tests = []struct {
		value    string
		expected string
	}{
		{
			"Starlink Ілона Маска відкриє офіс в Україні",
			"starlink-ilona-maska-vidkryie-ofis-v-ukraini",
		},
		{"Hellö Wörld", "hello-world"},
		{"你好世界", "ni-hao-shi-jie"},
		{"[^你好世界$]", "ni-hao-shi-jie"},
		{"This & that", "this-and-that"},
		{"\tHellö \t Wörld\n ", "hello-world"},
	}

	s := New()
	s.Lang(lang.UK)
	for _, test := range tests {
		if v := s.Lower(test.value); v != test.expected {
			t.Errorf("expected %s but %s", test.expected, v)
		}
	}
}
