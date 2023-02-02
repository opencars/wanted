package query

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/opencars/seedwork"
	"github.com/opencars/translit"
	"github.com/opencars/wanted/pkg/domain/model"
)

type Find struct {
	VINs    []string
	Numbers []string
}

func (q *Find) Prepare() {
	for i, number := range q.Numbers {
		q.Numbers[i] = translit.ToUA(strings.ToUpper(number))
	}

	for i, vin := range q.VINs {
		q.VINs[i] = translit.ToLatin(strings.ToUpper(vin))
	}
}

func (q *Find) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(
			&q.VINs,
			validation.Each(
				validation.Required.Error(seedwork.Required),
				validation.Length(6, 18).Error(seedwork.Invalid),
			),
		),
		validation.Field(
			&q.Numbers,
			validation.Each(
				validation.Required.Error(seedwork.Required),
				validation.Length(2, 18).Error(seedwork.Invalid),
			),
		),
	)

}

type FindResult struct {
	Vehicles []model.Vehicle
}
