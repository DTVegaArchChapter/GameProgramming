package wordle

import (
	"unicode"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	TurkishUpper = cases.Upper(language.Turkish)
)

func contains(runes []rune, r rune) bool {
	for _, v := range runes {
		if unicode.TurkishCase.ToUpper(v) == unicode.TurkishCase.ToUpper(r) {
			return true
		}
	}
	return false
}
