package cleansing

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/config"
)

var (
	settings *config.Cleansing
)

func TestCleansing_Brand_1(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"ВАЗ - 2107", "ВАЗ", "2107"},
		{"ВАЗ - 21099", "ВАЗ", "21099"},
		{"ВАЗ - 21063", "ВАЗ", "21063"},
		{"ВАЗ - 2101", "ВАЗ", "2101"},
		{"ВАЗ - 21093", "ВАЗ", "21093"},
		{"ВАЗ - 2106", "ВАЗ", "2106"},
		{"BAЗ21063", "ВАЗ", "21063"},
		{"BAЗ21099", "ВАЗ", "21099"},
		{"BAЗ2107", "ВАЗ", "2107"},
		{"ВАЗ - 21011", "ВАЗ", "21011"},
		{"BAЗ2101", "ВАЗ", "2101"},
		{"ВАЗ - 21013", "ВАЗ", "21013"},
		{"BAЗ21093", "ВАЗ", "21093"},
		{"BAЗ2106", "ВАЗ", "2106"},
		{"ВАЗ - 2103", "ВАЗ", "2103"},
		{"ВАЗ - 21061", "ВАЗ", "21061"},
		{"ВАЗ - 2105", "ВАЗ", "2105"},
		{"ВАЗ - 2108", "ВАЗ", "2108"},
		{"ВАЗ", "ВАЗ", ""},
		{"ВАЗ - 2109", "ВАЗ", "2109"},
	}

	for _, tt := range tests {
		maker, model, err := c.Brand(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.maker, maker)
		assert.Equal(t, tt.model, model)
	}
}

func TestCleansing_Brand_2(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"JAWA350", "JAWA", "350"},
		{"JAWA350638", "JAWA", "350638"},
		{"JAWA350634", "JAWA", "350634"},
		{"JAWA - 350", "JAWA", "350"},
		{"ЯBA350", "JAWA", "350"},
		{"ЯВА350", "JAWA", "350"},
		{"JAWA", "JAWA", ""},
		{"JAWA250", "JAWA", "250"},
		{"JAWA 350", "JAWA", "350"},
		{"ЯВА", "JAWA", ""},
		{"JAWA - 350 638", "JAWA", "350 638"},
	}

	for _, tt := range tests {
		maker, model, err := c.Brand(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.maker, maker)
		assert.Equal(t, tt.model, model)
	}
}

func TestCleansing_Brand_3(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"SUZUKI", "SUZUKI", ""},
		{"СУЗУКИ", "SUZUKI", ""},
		{"SUZUKI - LETS", "SUZUKI", "LETS"},
		{"CУЗУKИ", "SUZUKI", ""},
		{"SUZUKI - GRAND VITARA", "SUZUKI", "GRAND VITARA"},
		{"SUZUKI LETS", "SUZUKI", "LETS"},
		{"SUZUKI LETS 2", "SUZUKI", "LETS 2"},
		{"SUZUKI - SWIFT", "SUZUKI", "SWIFT"},
		{"SUZUKI SEPIA", "SUZUKI", "SEPIA"},
		{"SUZUKI - LETS 2", "SUZUKI", "LETS 2"},
		{"SUZUKI - SEPIA", "SUZUKI", "SEPIA"},
		{"SUZUKI ADDRESS", "SUZUKI", "ADDRESS"},
		{"SUZUKI - VERDE", "SUZUKI", "VERDE"},
		{"SUZUKISEPIA", "SUZUKI", "SEPIA"},
		{"SUZUKІ", "SUZUKI", ""},
		{"SUZUKILETS", "SUZUKI", "LETS"},
		{"SUZUKILETS2", "SUZUKI", "LETS2"},
		{"SUZUKI GRAND VITARA", "SUZUKI", "GRAND VITARA"},
		{"SUZUKI - VITARA", "SUZUKI", "VITARA"},
		{"SUZUKI ZZ", "SUZUKI", "ZZ"},
		{"SUZUKIADDRESS", "SUZUKI", "ADDRESS"},
		{"SUZUKI - BANDIT 250", "SUZUKI", "BANDIT 250"},
		{"СУЗУКІ", "SUZUKI", ""},
	}

	for _, tt := range tests {
		maker, model, err := c.Brand(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.maker, maker)
		assert.Equal(t, tt.model, model)
	}
}

func TestCleansing_Brand_4(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"HONDA", "HONDA", ""},
		{"HONDA - ACCORD", "HONDA", "ACCORD"},
		{"HONDA DIO", "HONDA", "DIO"},
		{"HONDA - DIO", "HONDA", "DIO"},
		{"ХОНДА", "HONDA", ""},
		{"HONDA - CR V", "HONDA", "CR-V"},
		{"HONDA - CR-V", "HONDA", "CR-V"},
		{"HONDA - CIVIC", "HONDA", "CIVIC"},
		{"HONDADIO", "HONDA", "DIO"},
		{"ХОНДА ДИО", "HONDA", "DIO"},
		{"HONDA-DIO", "HONDA", "DIO"},
		{"HONDA TACT", "HONDA", "TACT"},
		{"HONDA LEAD", "HONDA", "LEAD"},
		{"HONDA - LEAD", "HONDA", "LEAD"},
		{"HONDA - TACT", "HONDA", "TACT"},
		{"ХОНДАДИО", "HONDA", "DIO"},
		{"HONDA - CІVІC", "HONDA", "CІVІC"},
		{"HONDA - TAKT", "HONDA", "TAKT"},
		{"HONDA - PILOT", "HONDA", "PILOT"},
		{"HONDA - GIORNO", "HONDA", "GIORNO"},
		{"HONDA TAKT", "HONDA", "TAKT"},
		{"ХОНДА ТАКТ", "HONDA", "ТАКТ"},
		{"HONDA CR-V", "HONDA", "CR-V"},
		{"ХОНДАТАКТ", "HONDA", "ТАКТ"},
		{"HONDA DIO 27", "HONDA", "DIO 27"},
		{"HONDALEAD", "HONDA", "LEAD"},
		{"HONDA - ACCORD 2.0", "HONDA", "ACCORD 2.0"},
		{"HONDATACT", "HONDA", "TACT"},
		{"HONDAACCORD", "HONDA", "ACCORD"},
		{"HONDA - DIO AF34", "HONDA", "DIO AF34"},
		{"HONDATAKT", "HONDA", "TAKT"},
	}

	for _, tt := range tests {
		maker, model, err := c.Brand(tt.in)
		assert.NoError(t, err)
		assert.Equal(t, tt.maker, maker)
		assert.Equal(t, tt.model, model)
	}
}

func TestCleansing_Brand_5(t *testing.T) {
	c := New(settings)

	var tests = []struct {
		in           string
		maker, model string
	}{
		{"DAEWOO - LANOS", "DAEWOO", "LANOS"},
		{"DAEWOO", "DAEWOO", ""},
		{"DAEWOO - SENS", "DAEWOO", "SENS"},
		{"DAEWOOLANOS", "DAEWOO", "LANOS"},
		{"DAEWOO - NEXIA", "DAEWOO", "NEXIA"},
		{"DAEWOO - LANOS 1.5", "DAEWOO", "LANOS 1.5"},
		{"DAEWOO - MATIZ", "DAEWOO", "MATIZ"},
		{"DAEWOO - LANOS TF69Y", "DAEWOO", "LANOS TF69Y"},
		{"DAEWOO LANOS TF69Y", "DAEWOO", "LANOS TF69Y"},
		{"DAEWOO LANOS", "DAEWOO", "LANOS"},
		{"DАЕWОО-LАNОS", "DAEWOO", "LANOS"},
		{"DAEWOOSENS", "DAEWOO", "SENS"},
		{"DAEWOO - MATІZ", "DAEWOO", "MATІZ"},
		{"DAEWOO - NUBIRA", "DAEWOO", "NUBIRA"},
		{"DAEWOO - ESPERO", "DAEWOO", "ESPERO"},
		{"DAEWOO - LANOS D4LM500", "DAEWOO", "LANOS D4LM500"},
		{"DAEWOOLANOSTF69Y", "DAEWOO", "LANOS TF69Y"},
		{"DAEWOO - LANOS T13110", "DAEWOO", "LANOS T13110"},
		{"DAEWOO-LANOS", "DAEWOO", "LANOS"},
		{"DEO MATIZ", "DAEWOO", "MATIZ"},
		{"ДЭО ЛАНОС", "DAEWOO", "LANOS"},
		{"ДЭО СЕНС", "DAEWOO", "SENS"},
		{"ДЭОСЕНС", "DAEWOO", "SENS"},
		{"ДЭОЛАНОС", "DAEWOO", "LANOS"},
		{"ДЭУ", "DAEWOO", ""},
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
