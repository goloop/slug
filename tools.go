package slug

import (
	"strings"

	"github.com/goloop/t13n/lang"
)

// The slugRules sets custom transliteration rules
// for the github.com/goloop/t13n module.
func slugRules(ts lang.TransState) (string, int, bool) {
	// Ignore ranges.
	// Important: edit maps ascending order only!
	ignored := [][]int{
		{0, 47},
		{58, 64},
		{91, 96},
		{123, 141},
		{143, 152},
		{155, 155},
	}

	switch ts.Value {
	case " ", "~", "_", "\t", "\n":
		ts.Value = "-"
	case "@":
		ts.Value = "at"
	case "&":
		ts.Value = "and"
	case "#":
		ts.Value = "sharp"
	case "%":
		ts.Value = "pct"
	default:
		id := int(ts.Curr)
		for _, d := range ignored {
			// If the item isn't in the following ranges.
			if id < d[0] {
				break
			}

			// If the item is in the current range.
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
