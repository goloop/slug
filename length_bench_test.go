package slug

import (
	"strings"
	"testing"
)

var limitedResult string

func BenchmarkLimitedLongSlug(b *testing.B) {
	s := New(WithMaxLength(32))
	text := strings.Repeat("long text to shorten ", 1000)
	s.Make(text)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limitedResult = s.Make(text)
	}
}
