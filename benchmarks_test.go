package slug

import (
	"testing"
)

var (
	testString   = "Hello 世界 @ #$% &~_ Testing-Slug"
	resultString string
	benchResults = make([]string, 0, 100)
)

func BenchmarkMake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultString = Make(testString)
	}
}

func BenchmarkLower(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultString = Lower(testString)
	}
}

func BenchmarkUpper(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultString = Upper(testString)
	}
}

func BenchmarkSlugMake(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		resultString = s.Make(testString)
	}
}

func BenchmarkSlugLower(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		resultString = s.Lower(testString)
	}
}

func BenchmarkSlugUpper(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		resultString = s.Upper(testString)
	}
}

// Benchmark з різними мовними налаштуваннями
func BenchmarkSlugWithLang(b *testing.B) {
	s := New()
	langs := []string{"uk", "en", "de", "fr"}

	for _, l := range langs {
		b.Run("Lang_"+l, func(b *testing.B) {
			s.Lang(l)
			for i := 0; i < b.N; i++ {
				resultString = s.Make(testString)
			}
		})
	}
}
