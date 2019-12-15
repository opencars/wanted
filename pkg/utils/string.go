package utils

import "strings"

// Trim removes trailing spaces and symbols from the string.
func Trim(lexeme *string) *string {
	if lexeme == nil {
		return nil
	}

	str := *lexeme
	str = strings.TrimFunc(str, func(r rune) bool {
		return r == '-' || r == '%' || r == '*' || r == '.' || r == ' '
	})

	if str == "" {
		return nil
	}

	return &str
}
