package storage

type Database interface {
	WantedVehiclesByNumber(number string) ([]WantedVehicle, error)
	WantedVehiclesByVIN(vin string) ([]WantedVehicle, error)
}

type Store struct {
	db Database
}

// New returns new storage.
func New(db Database) *Store {
	return &Store{
		db: db,
	}
}

// WantedVehiclesByNumber returns all wanted vehicles by number.
func (s *Store) WantedVehiclesByNumber(number string) ([]WantedVehicle, error) {
	return s.db.WantedVehiclesByNumber(number)
}

// WantedVehiclesByVIN returns all wanted vehicles by vin.
func (s *Store) WantedVehiclesByVIN(vin string) ([]WantedVehicle, error) {
	return s.db.WantedVehiclesByVIN(vin)
}
