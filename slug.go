package slug

import (
	"regexp"
	"strings"

	"github.com/goloop/slug/lang"
	"github.com/goloop/t13n"
)

var (
	// Characters that are not allowed in the slug.
	wrongRegx = regexp.MustCompile(`([^A-Za-z0-9\-]|-$|^-)`)

	// Duplicate separator.
	dusepRegx = regexp.MustCompile(`-[-]+`)
)

// New retursn pointer to Slug.
func New() *Slug {
	return &Slug{lang: lang.None}
}

// Slug is the slug constructor.
type Slug struct {
	lang string
}

// Lang sets the type of language features to use during slugify.
func (s *Slug) Lang(l string) *Slug {
	s.lang = l
	return s
}

// Make returns slug from string.
func (s *Slug) Make(t string) string {
	var result string

	result = t13n.Render(s.lang, t, slugRules)
	result = dusepRegx.ReplaceAllString(result, "-")
	result = wrongRegx.ReplaceAllString(result, "")

	return result
}

// Lower returns slug in lowercase.
func (s *Slug) Lower(t string) string {
	return strings.ToLower(s.Make(t))
}

// Upper returns slug in uppercase.
func (s *Slug) Upper(t string) string {
	return strings.ToUpper(s.Make(t))
}
