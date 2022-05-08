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
