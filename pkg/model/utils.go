package model

import (
	"strings"
	"unicode"
)

func ParseBrand(lexeme string) (other string) {
	other = strings.TrimFunc(lexeme, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	return other
}
