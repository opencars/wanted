package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/opencars/wanted/pkg/cleansing"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/store"
)

type Store struct {
	db                 *sqlx.DB
	clean              *cleansing.Cleansing
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

func New(conf *config.Settings) (*Store, error) {
	info := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",	
		conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Name, conf.DB.SSLMode, conf.DB.Password,
	)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &Store{
		db:    db,
		clean: cleansing.New(&conf.Cleansing),
	}, nil
}
