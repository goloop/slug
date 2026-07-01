package slug

import (
	"strings"
	"testing"

	"github.com/goloop/slug/v2/lang"
)

var (
	testString   = "Hello 世界 @ #$% &~_ Testing-Slug"
	asciiString  = "The Quick Brown Fox Jumps Over The Lazy Dog"
	longString   = strings.Repeat("Привіт, світ! Hello 世界. ", 40) // > 256 runes
	resultString string
)

func BenchmarkMake(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultString = Make(testString)
	}
}

func BenchmarkLower(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultString = Lower(testString)
	}
}

func BenchmarkUpper(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultString = Upper(testString)
	}
}

// BenchmarkMakeASCII measures the common case of already-Latin input.
func BenchmarkMakeASCII(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultString = Make(asciiString)
	}
}

// BenchmarkMakeLong exercises the multi-chunk transliteration path.
func BenchmarkMakeLong(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resultString = Make(longString)
	}
}

func BenchmarkSlugWithLang(b *testing.B) {
	for _, code := range []string{lang.UK, lang.EN, lang.DE, lang.FR} {
		s := New(WithLang(code))
		b.Run("Lang_"+code, func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				resultString = s.Make(testString)
			}
		})
	}
}
