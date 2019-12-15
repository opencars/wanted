package store

import (
	"github.com/opencars/wanted/pkg/model"
)

type RevisionRepository interface {
	Create(revision *model.Revision) error
	FindByID(id string) (*model.Revision, error)
	Last() (*model.Revision, error)
	All() ([]model.Revision, error)
	AllIDs() ([]string, error)
	Stats() ([]model.RevisionStatMonth, error)
}

type VehicleRepository interface {
	Create(revision *model.Revision, added []model.Vehicle, removed []string) error
	FindByNumber(number string) ([]model.Vehicle, error)
	FindByVIN(vin string) ([]model.Vehicle, error)
	FindByRevisionID(id string) ([]model.Vehicle, error)
	All() ([]model.Vehicle, error)
	AllWithLimit(limit int64) ([]model.Vehicle, error)
}
