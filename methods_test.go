package slug

import (
	"strings"
	"testing"
)

// TestVersion tests the package version.
// Note: each time you change the major version, you need to fix the tests.
func TestVersion(t *testing.T) {
	var expected = "v0." // change it for major version

	version := Version()
	if strings.HasPrefix(version, expected) != true {
		t.Error("incorrect version")
	}

	if len(strings.Split(version, ".")) != 3 {
		t.Error("version format should be as " +
			"v{major_version}.{minor_version}.{patch_version}")
	}
}

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
		value    string
		expected string
	}{
		{
			"Starlink Ілона Маска відкриє офіс в Україні",
			"starlink-ilona-maska-vidkriie-ofis-v-ukrayini",
		},
		{"Hellö Wörld", "hello-world"},
		{"你好世界", "ni-hao-shi-jie"},
		{"[^你好世界$]", "ni-hao-shi-jie"},
		{"This & that", "this-and-that"},
		{"\tHellö \t Wörld\n ", "hello-world"},
	}

	for _, test := range tests {
		if v := Lower(test.value); v != test.expected {
			t.Errorf("expected %s but %s", test.expected, v)
		}
	}
}
