package teststore

import (
	"github.com/opencars/wanted/pkg/model"
)

type VehicleRepository struct {
	store    *Store
	vehicles map[string]*model.Vehicle
}

func (r *VehicleRepository) Create(revision *model.Revision, vehicles ...model.Vehicle) error {
	return nil
}

func (r *VehicleRepository) All() ([]model.Vehicle, error) {
	return nil, nil
}

func (r *VehicleRepository) CreateRevisionAndAll(revision *model.Revision, vehicles []model.WantedVehicle) error {
	return nil
}

func (r *VehicleRepository) FindByNumber(number string) ([]model.Vehicle, error) {
	return nil, nil
}

func (r *VehicleRepository) FindByVIN(vin string) ([]model.Vehicle, error) {
	return nil, nil
}

func (r *VehicleRepository) FindByRevisionID(id string) ([]model.Vehicle, error) {
	return nil, nil
}

func (r *VehicleRepository) AllWithLimit(limit int64) ([]model.Vehicle, error) {
	return nil, nil
}
