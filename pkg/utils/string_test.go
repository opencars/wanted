package utils_test

import (
	"fmt"
	"testing"

	"github.com/opencars/wanted/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestTrim(t *testing.T) {
	t.Run("removes bad characters", func(t *testing.T) {
		tests := []struct {
			in, out string
		}{
			{" Hyundai ", "Hyundai"},
			{" BMW .* ", "BMW"},
		}

		for _, test := range tests {
			assert.Equal(t, test.out, *utils.Trim(&test.in))
		}
	})

	t.Run("returns nil", func(t *testing.T) {
		tests := []string{
			"   ", "***", "---", "%%%", "...",
			"*-%.%-*", " *-. . .-* ",
		}

		for _, test := range tests {
			assert.Nil(t, utils.Trim(&test))
		}
	})
}

func ExampleTrim() {
	input := "... this string will be fixed%. "
	fmt.Println(*utils.Trim(&input))
	// Output: this string will be fixed
}
