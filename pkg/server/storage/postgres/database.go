package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/opencars/wanted/pkg/server/storage"
)

type database struct {
	db *sqlx.DB
}

func New(host string, port int, user, password, dbname string) (storage.Database, error) {
	info := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sqlx.Connect("postgres", info)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	return &database{
		db: db,
	}, nil
}

func (db *database) WantedVehiclesByNumber(number string) ([]storage.WantedVehicle, error) {
	vehicles := make([]storage.WantedVehicle, 0)

	err := db.db.Select(&vehicles, `SELECT * FROM wanted_vehicles WHERE number = $1`, number)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}

func (db *database) WantedVehiclesByVIN(vin string) ([]storage.WantedVehicle, error) {
	vehicles := make([]storage.WantedVehicle, 0)

	err := db.db.Select(&vehicles, `SELECT * FROM wanted_vehicles WHERE body_number = $1 OR chassis_number = $1 OR engine_number = $1`, vin)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}
