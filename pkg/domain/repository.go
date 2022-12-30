package domain

import (
	"context"

	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/domain/query"
)

type RevisionRepository interface {
	Create(revision *model.Revision) error
	FindByID(id string) (*model.Revision, error)
	Last() (*model.Revision, error)
	All() ([]model.Revision, error)
	AllWithLimit(limit uint64) ([]model.Revision, error)
	AllIDs() ([]string, error)
	Stats() ([]model.RevisionStatMonth, error)
}

type VehicleRepository interface {
	Create(revision *model.Revision, added []model.Vehicle, removed []string) error
	Find(context.Context, *query.Find) (*query.FindResult, error)
	FindByNumber(number string) ([]model.Vehicle, error)
	FindByVIN(vin string) ([]model.Vehicle, error)
	FindByRevisionID(id string) ([]model.Vehicle, error)
	All() ([]model.Vehicle, error)
	AllWithLimit(limit uint64) ([]model.Vehicle, error)
}

type Store interface {
	Revision() RevisionRepository
	Vehicle() VehicleRepository
}
