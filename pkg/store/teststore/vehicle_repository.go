package teststore

import (
	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store"
)

type VehicleRepository struct {
	store    *Store
	vehicles map[string]*model.Vehicle
}

func (r *VehicleRepository) Create(revision *model.Revision, added []model.Vehicle, removed []string) error {
	if err := r.store.Revision().Create(revision); err != nil {
		return err
	}

	for i, v := range added {
		if _, ok := r.vehicles[v.ID]; !ok {
			r.vehicles[v.ID] = &added[i]
		} else {
			r.vehicles[v.ID].Status = added[i].Status
		}
	}

	for _, id := range removed {
		if _, ok := r.vehicles[id]; !ok {
			return store.ErrRecordNotFound
		} else {
			r.vehicles[id].Status = model.StatusRemoved
		}
	}

	return nil
}

func (r *VehicleRepository) All() ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0, len(r.vehicles))

	for _, v := range r.vehicles {
		vehicles = append(vehicles, *v)
	}

	return vehicles, nil
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
