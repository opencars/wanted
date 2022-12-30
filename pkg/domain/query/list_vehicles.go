package query

import (
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/opencars/seedwork"
	"github.com/opencars/wanted/pkg/domain/model"
)

type ListVehicles struct {
	UserID  string
	TokenID string
	Limit   string
	Offset  string
}

func (q *ListVehicles) Prepare() {}

func (q *ListVehicles) GetLimit() uint64 {
	if q.Limit == "" {
		return 10
	}

	num, err := strconv.ParseInt(q.Limit, 10, 64)
	if err != nil {
		panic(err)
	}

	if num < 0 {
		return 10
	}

	return uint64(num)
}

func (q *ListVehicles) GetOffset() uint64 {
	if q.Offset == "" {
		return 0
	}

	num, err := strconv.ParseInt(q.Offset, 10, 64)
	if err != nil {
		panic(err)
	}

	if num < 0 {
		return 10
	}

	return uint64(num)
}

func (q *ListVehicles) Validate() error {
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
			&q.Limit,
			is.Int.Error(seedwork.IsNotInreger),
		),
		validation.Field(
			&q.Offset,
			is.Int.Error(seedwork.IsNotInreger),
		),
	)
}

type ListVehiclesResult struct {
	Vehicles []model.Vehicle
}
