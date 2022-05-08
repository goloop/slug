package slug

import (
	"testing"
)

// TestMake tests Make function.
func TestMake(t *testing.T) {
	var tests = []struct {
		value    string
		expected string
	}{
		{
			"Starlink Ілона Маска відкриє офіс в Україні",
			"Starlink-Ilona-Maska-vidkriie-ofis-v-Ukrayini",
		},
		{"Hellö Wörld", "Hello-World"},
		{"你好世界", "Ni-Hao-Shi-Jie"},
		{"[^你好世界$]", "Ni-Hao-Shi-Jie"},
		{"This & that", "This-and-that"},
		{"\tHellö \t Wörld\n ", "Hello-World"},
	}

	for _, test := range tests {
		if v := Make(test.value); v != test.expected {
			t.Errorf("expected %s but %s", test.expected, v)
		}
	}
}

// TestLower tests Lower function.
func TestLower(t *testing.T) {
	var tests = []struct {
		value         string
		lowerExpected string
		upperExpected string
	}{
		{
			"Starlink Ілона Маска відкриє офіс в Україні",
			"starlink-ilona-maska-vidkriie-ofis-v-ukrayini",
			"STARLINK-ILONA-MASKA-VIDKRIIE-OFIS-V-UKRAYINI",
		},
		{"Hellö Wörld", "hello-world", "HELLO-WORLD"},
		{"你好世界", "ni-hao-shi-jie", "NI-HAO-SHI-JIE"},
		{"[^你好世界$]", "ni-hao-shi-jie", "NI-HAO-SHI-JIE"},
		{"This & that", "this-and-that", "THIS-AND-THAT"},
		{"\tHellö \t Wörld\n ", "hello-world", "HELLO-WORLD"},
	}

	for _, test := range tests {
		if v := Lower(test.value); v != test.lowerExpected {
			t.Errorf("expected %s but %s", test.lowerExpected, v)
		}

		if v := Upper(test.value); v != test.upperExpected {
			t.Errorf("expected %s but %s", test.upperExpected, v)
		}
	}
}
