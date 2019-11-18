package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: {input: "ГАЗ (Легковий))", other: "ГАЗ", kind: "Легковий"},
// TODO: {input: "Volkswagen Polo", other: "Volkswagen Polo", kind: ""},
func TestParseKind(t *testing.T) {
	type test struct {
		input       string
		other, kind string
	}

	tests := []test{
		{input: "Lexus - RX300 (Приклад)", other: "Lexus - RX300", kind: "Приклад"},
		{input: "Hyundai - i30 (Легковий автотранспорт)", other: "Hyundai - i30", kind: "Легковий автотранспорт"},
		{input: "KOGEL (Автопричіп)", other: "KOGEL", kind: "Автопричіп"},
		{input: "BMW - i7 ( )", other: "BMW - i7", kind: ""},
		{input: "(Автопричіп)", other: "", kind: "Автопричіп"},
		{input: "ВАЗ ((Легковий)", other: "ВАЗ", kind: "Легковий"},
		{input: "Skoda - Octavia ((Легковий)))", other: "Skoda - Octavia", kind: "Легковий"},
	}

	for _, tc := range tests {
		other, kind := ParseKind(tc.input)
		assert.Equal(t, tc.other, other, "other")
		assert.Equal(t, tc.kind, kind, "kind")
	}
}
