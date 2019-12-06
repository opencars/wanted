package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/opencars/wanted/pkg/store"
)

type Store struct {
	db                 *sqlx.DB
	revisionRepository *RevisionRepository
	vehicleRepository  *VehicleRepository
}

func (s *Store) Revision() store.RevisionRepository {
	if s.revisionRepository != nil {
		return s.revisionRepository
	}

	s.revisionRepository = &RevisionRepository{
		store: s,
	}

	return s.revisionRepository
}

func (s *Store) Vehicle() store.VehicleRepository {
	if s.vehicleRepository != nil {
		return s.vehicleRepository
	}

	s.vehicleRepository = &VehicleRepository{
		store: s,
	}

	return s.vehicleRepository
}

func New(host string, port int, user, password, dbname string) (*Store, error) {
	info := fmt.Sprintf("host=%s port=%d user=%s password=\"%s\" dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Store{
		db: db,
	}, nil
}
