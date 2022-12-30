package query

import (
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/opencars/schema"
	"github.com/opencars/seedwork"
	"github.com/opencars/translit"

	"github.com/opencars/wanted/pkg/domain/model"
)

type ListByNumber struct {
	UserID  string
	TokenID string
	Number  string
}

func (q *ListByNumber) Prepare() {
	q.Number = translit.ToUA(strings.ToUpper(q.Number))
}

func (q *ListByNumber) Validate() error {
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
			&q.Number,
			validation.Required.Error(seedwork.Required),
			validation.Length(6, 18).Error(seedwork.Invalid),
		),
	)
}

func (q *ListByNumber) Event(operations ...model.Vehicle) schema.Producable {
	// msg := vehicle.OperationSearched{
	// 	UserId:       q.UserID,
	// 	TokenId:      q.TokenID,
	// 	Number:       q.Number,
	// 	ResultAmount: uint32(len(operations)),
	// 	SearchedAt:   timestamppb.New(time.Now().UTC()),
	// }

	// return schema.New(&source, &msg).Message(
	// 	schema.WithSubject(schema.OperationCustomerActions),
	// )
	return nil
}
