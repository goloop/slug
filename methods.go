package slug

import (
	"strings"

	"github.com/goloop/t13n"
	"github.com/goloop/t13n/lang"
)

// Make returns slug from string.
func Make(t string) string {
	var result string

	result = t13n.Render(lang.None, t, slugRules)
	result = dusepRegx.ReplaceAllString(result, "-")
	result = wrongRegx.ReplaceAllString(result, "")

	return result
}

// Lower returns slug in lowercase.
func Lower(t string) string {
	return strings.ToLower(Make(t))
}

// Upper returns slug in uppercase.
func Upper(t string) string {
	return strings.ToUpper(Make(t))
}
