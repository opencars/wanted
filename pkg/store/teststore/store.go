package teststore

import (
	"github.com/opencars/wanted/pkg/model"

	"github.com/opencars/wanted/pkg/store"
)

type Store struct {
	revisionRepository *RevisionRepository
	vehicleRepository  *VehicleRepository
}

func (s *Store) Revision() store.RevisionRepository {
	if s.revisionRepository != nil {
		return s.revisionRepository
	}

	s.revisionRepository = &RevisionRepository{
		store:     s,
		revisions: make(map[string]*model.Revision),
	}

	return s.revisionRepository
}

func (s *Store) Vehicle() store.VehicleRepository {
	if s.vehicleRepository != nil {
		return s.vehicleRepository
	}

	s.vehicleRepository = &VehicleRepository{
		store:    s,
		vehicles: make(map[string]*model.Vehicle),
	}

	return s.vehicleRepository
}

func New() *Store {
	return &Store{}
}
