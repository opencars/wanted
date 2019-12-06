package model

import (
	"strings"
	"unicode"
)

// ParseKind splits lexeme into kind and other lexemes.
//
// Example:
// ParseKind("Tesla - Model X (Легковий автотранспорт)")
//  =>
// other: "Hyundai - i30", kind: "Легковий автотранспорт"
func ParseKind(lexeme string) (other string, kind string) {
	stack := make([]rune, 0)

	for i := len(lexeme) - 1; i >= 0; i-- {
		switch lexeme[i] {
		case ')':
			stack = append(stack, ')')
		case '(':
			stack = stack[:len(stack)-1]
		}

		if len(stack) == 0 {
			other, kind = lexeme[:i], lexeme[i:]
			break
		}
	}

	if len(stack) != 0 {
		var tmp []rune
		i := 0
		for ; i < len([]rune(lexeme)) && lexeme[i] != '('; i++ {
			tmp = append(tmp, []rune(lexeme)[i])
		}
		other = string(tmp)
		kind = lexeme[i:]
	}

	other = strings.TrimFunc(other, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r)
	})

	kind = strings.TrimFunc(kind, func(r rune) bool {
		return !unicode.IsLetter(r)
	})

	return other, kind
}
