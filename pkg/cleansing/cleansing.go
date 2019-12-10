package cleansing

import (
	"errors"
	"regexp"

	"github.com/opencars/wanted/pkg/config"
)

type Cleansing struct {
	conf *config.Cleansing
}

func New(conf *config.Cleansing) *Cleansing {
	return &Cleansing{
		conf: conf,
	}
}

func (c *Cleansing) Brand(lexeme string) (string, string, error) {
	for _, m := range c.conf.Brand.Matchers {
		r, err := regexp.Compile(m.Pattern)
		if err != nil {
			return "", "", err
		}

		if !r.MatchString(lexeme) {
			continue
		}

		maker := r.ReplaceAllString(lexeme, m.Maker)
		model := r.ReplaceAllString(lexeme, m.Model)

		return maker, model, nil
	}

	return "", "", errors.New("failed to match")
}
