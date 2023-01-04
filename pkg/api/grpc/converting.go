package grpc

import (
	"github.com/araddon/dateparse"
	"github.com/opencars/grpc/pkg/common"
	"github.com/opencars/grpc/pkg/wanted"

	"github.com/opencars/wanted/pkg/domain/model"
)

func fromDomain(dto *model.Vehicle) *wanted.Vehicle {
	var vehicle wanted.Vehicle

	vehicle.Id = dto.CheckSum
	vehicle.RevisionId = dto.RevisionID
	vehicle.Ovd = dto.OVD

	if dto.Brand != nil {
		vehicle.Title = *dto.Brand
	}

	if dto.Model != nil {
		vehicle.Model = *dto.Model
	}

	if dto.Maker != nil {
		vehicle.Brand = *dto.Maker
	}

	if dto.Color != nil {
		vehicle.Color = *dto.Color
	}

	if dto.Number != nil {
		vehicle.Number = *dto.Number
	}

	if dto.BodyNumber != nil {
		vehicle.BodyNumber = *dto.BodyNumber
	}

	if dto.ChassisNumber != nil {
		vehicle.ChassisNumber = *dto.ChassisNumber
	}

	if dto.EngineNumber != nil {
		vehicle.EngineNumber = *dto.EngineNumber
	}

	if dto.Kind != nil {
		vehicle.Kind = *dto.Kind
	}

	switch dto.Status {
	case model.StatusStolen:
		vehicle.Status = wanted.Vehicle_STATUS_STOLEN
	case model.StatusRemoved:
		vehicle.Status = wanted.Vehicle_STATUS_REMOVED
	default:
		vehicle.Status = wanted.Vehicle_STATUS_UNKNOWN
	}

	t, _ := dateparse.ParseAny(dto.TheftDate)

	vehicle.TheftDate = &common.Date{
		Year:  int32(t.Year()),
		Month: int32(t.Month()),
		Day:   int32(t.Day()),
	}

	vehicle.InsertDate = &common.Date{
		Year:  int32(dto.InsertDate.Year()),
		Month: int32(dto.InsertDate.Month()),
		Day:   int32(dto.InsertDate.Day()),
	}

	return &vehicle
}
