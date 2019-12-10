package cleansing

import (
	"os"
	"testing"

	"github.com/opencars/wanted/pkg/config"
	"github.com/stretchr/testify/assert"
)

var (
	settings *config.Cleansing
)

func TestCleansing_Brand(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"ВАЗ 210994", "ВАЗ", "210994"},
		{"ВАЗ - 210994", "ВАЗ", "210994"},
		{"ВАЗ210994", "ВАЗ", "210994"},
		{"ВАЗ -210994", "ВАЗ", "210994"},
		{"ВА3  210994", "ВАЗ", "210994"},
		{"ВAЗ   210994", "ВАЗ", "210994"},
		{"BАЗ - 210994", "ВАЗ", "210994"},
		{"BA3-XYZ", "ВАЗ", "XYZ"},
	}

	for _, tt := range tests {
		maker, model, err := c.Brand(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.maker, maker)
		assert.Equal(t, tt.model, model)
	}
}

func TestMain(m *testing.M) {
	conf, err := config.New("../../config/config.toml")
	if err != nil {
		panic(err)
	}

	settings = &conf.Cleansing
	os.Exit(m.Run())
}
