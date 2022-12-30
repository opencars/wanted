package teststore

import (
	"context"

	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/domain/query"
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
		if _, ok := r.vehicles[v.CheckSum]; !ok {
			r.vehicles[v.CheckSum] = &added[i]
		} else {
			r.vehicles[v.CheckSum].Status = added[i].Status
		}
	}

	for _, id := range removed {
		if _, ok := r.vehicles[id]; !ok {
			return model.ErrNotFound
		}

		r.vehicles[id].Status = model.StatusRemoved
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

func (r *VehicleRepository) AllWithLimit(limit uint64) ([]model.Vehicle, error) {
	return nil, nil
}

func (r *VehicleRepository) Find(context.Context, *query.Find) (*query.FindResult, error) {
	return nil, nil
}
