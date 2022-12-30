package query

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/opencars/schema"
	"github.com/opencars/seedwork"
	"github.com/opencars/translit"

	"github.com/opencars/wanted/pkg/domain/model"
)

type ListByVIN struct {
	UserID  string
	TokenID string
	VIN     string
}

func (q *ListByVIN) Prepare() {
	q.VIN = translit.ToLatin(strings.ToUpper(q.VIN))
}

func (q *ListByVIN) Validate() error {
	return validation.ValidateStruct(q,
		validation.Field(
			&q.UserID,
			validation.Required.Error(seedwork.Required),
		),
		validation.Field(
			&q.TokenID,
			validation.Required.Error(seedwork.Required),
		),
		validation.Field(
			&q.VIN,
			validation.Required.Error(seedwork.Required),
			validation.Length(6, 18).Error(seedwork.Invalid),
		),
	)
}

func (q *ListByVIN) Event(operations ...model.Vehicle) schema.Producable {
	// msg := vehicle.OperationSearched{
	// 	UserId:       q.UserID,
	// 	TokenId:      q.TokenID,
	// 	Vin:          q.VIN,
	// 	ResultAmount: uint32(len(operations)),
	// 	SearchedAt:   timestamppb.New(time.Now().UTC()),
	// }

	// return schema.New(&source, &msg).Message(
	// 	schema.WithSubject(schema.OperationCustomerActions),
	// )

	return nil
}
