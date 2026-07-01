package slug_test

import (
	"fmt"

	"github.com/goloop/slug/v2"
	"github.com/goloop/slug/v2/lang"
)

func ExampleMake() {
	fmt.Println(slug.Make("Hello World"))
	// Output: Hello-World
}

func ExampleLower() {
	fmt.Println(slug.Lower("Hello World"))
	// Output: hello-world
}

func ExampleUpper() {
	fmt.Println(slug.Upper("Hello World"))
	// Output: HELLO-WORLD
}

// Meaningful symbols become words and punctuation becomes a separator.
func ExampleMake_symbols() {
	fmt.Println(slug.Make("AT&T: 100% ready"))
	// Output: AT-and-T-100-pct-ready
}

// A hyphen in the input is preserved as a separator, not dropped.
func ExampleMake_hyphen() {
	fmt.Println(slug.Make("co-operate with e-mail"))
	// Output: co-operate-with-e-mail
}

func ExampleNew() {
	s := slug.New(slug.WithLang(lang.UK))
	fmt.Println(s.Make("Привіт, світ!"))
	// Output: Pryvit-svit
}

func ExampleWithSeparator() {
	s := slug.New(slug.WithSeparator("_"))
	fmt.Println(s.Make("Hello World"))
	// Output: Hello_World
}

func ExampleWithMaxLength() {
	s := slug.New(slug.WithMaxLength(11))
	fmt.Println(s.Make("hello world foo bar"))
	// Output: hello-world
}

func ExampleWithFallback() {
	s := slug.New(slug.WithFallback("post"))
	fmt.Println(s.Make("!!!"))
	// Output: post
}

func ExampleSlug_MakeUnique() {
	taken := map[string]bool{"hello-world": true, "hello-world-2": true}
	s := slug.New()
	fmt.Println(s.MakeUnique("hello world", func(x string) bool {
		return taken[x]
	}))
	// Output: hello-world-3
}
