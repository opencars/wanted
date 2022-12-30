package domain

import (
	"context"

	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/domain/query"
)

type InternalService interface {
	Find(context.Context, *query.Find) (*query.FindResult, error)
}

type CustomerService interface {
	FindRevisionByID(context.Context, string) (*model.Revision, error)
	ListRevisions(context.Context, *query.ListRevisions) ([]model.Revision, error)

	ListByNumber(context.Context, *query.ListByNumber) ([]model.Vehicle, error)
	ListByVIN(context.Context, *query.ListByVIN) ([]model.Vehicle, error)
	ListVehicles(context.Context, *query.ListVehicles) ([]model.Vehicle, error)
}
