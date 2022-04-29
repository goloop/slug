package slug

import (
	"strings"

	"github.com/goloop/t13n/lang"
)

// The slugRules sets additional transliteration rules
// for the github.com/goloop/t13n module.
func slugRules(ts lang.TransState) (string, int, bool) {
	var lib = map[string]string{
		" ":  "-",
		"~":  "-",
		"_":  "-",
		"\t": "-",
		"\n": "-",
		"@":  "at",
		"&":  "and",
		"#":  "harp",
		"%":  "percentage",
	}

	if v, ok := lib[ts.Value]; ok {
		ts.Value = v
	} else {
		var ignore = [][]int{
			{0, 47},
			{58, 64},
			{91, 96},
			{123, 141},
			{143, 152},
			{155, 155},
		}

		id := int(ts.Curr)
		for _, d := range ignore {
			if id < d[0] {
				break
			}

			if id >= d[0] && id <= d[1] {
				ts.Value = ""
				break
			}
		}
	}

	if len(ts.Value) > 1 && strings.HasSuffix(ts.Value, " ") {
		runes := []rune(ts.Value)
		if ts.Next != 0 {
			ts.Value = string(runes[:len(runes)-1]) + "-"
		}
	}

	return ts.Value, 0, true
}
