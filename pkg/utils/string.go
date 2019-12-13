package utils

import "strings"

func Trim(lexeme *string) *string {
	if lexeme == nil {
		return nil
	}

	str := *lexeme
	str = strings.TrimSpace(str)

	str = strings.TrimFunc(str, func(r rune) bool {
		return r == '-' || r == '%' || r == '*' || r == '.'
	})

	if str == "" {
		return nil
	}

	return &str
}
