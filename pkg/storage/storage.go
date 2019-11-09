package storage

type Database interface {
	CreateOrUpdateVehicle(*WantedVehicle) error
	AllVehicles() ([]WantedVehicle, error)
	AllRevisions() ([]Revision, error)
	AllRevisionIDs() ([]string, error)
	VehiclesByNumber(number string) ([]Vehicle, error)
	VehiclesByVIN(vin string) ([]Vehicle, error)
	Vehicles(limit int64) ([]Vehicle, error)

	CreateOrUpdateVehicles(vv []WantedVehicle) error
	CreateOrUpdateVehicles2(vv []WantedVehicle) error
	CreateRevisionAndVehicles(r *Revision, vv []WantedVehicle) error

	CreateRevision(*Revision) error
	Revisions(limit int64) ([]Revision, error)
	RevisionByID(id string) (*Revision, error)
}

type Store struct {
	db Database
}

func New(db Database) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) CreateOrUpdateVehicle(v *WantedVehicle) error {
	return s.db.CreateOrUpdateVehicle(v)
}

func (s *Store) CreateRevision(r *Revision) error {
	return s.db.CreateRevision(r)
}

func (s *Store) CreateRevisionAndVehicles(r *Revision, vv []WantedVehicle) error {
	return s.db.CreateRevisionAndVehicles(r, vv)
}

func (s *Store) AllVehicles() ([]WantedVehicle, error) {
	return s.db.AllVehicles()
}

func (s *Store) AllRevisions() ([]Revision, error) {
	return s.db.AllRevisions()
}

func (s *Store) AllRevisionIDs() ([]string, error) {
	return s.db.AllRevisionIDs()
}

func (s *Store) CreateOrUpdateVehicles(vv []WantedVehicle) error {
	return s.db.CreateOrUpdateVehicles(vv)
}

func (s *Store) CreateOrUpdateVehicles2(vv []WantedVehicle) error {
	return s.db.CreateOrUpdateVehicles2(vv)
}

// WantedVehiclesByNumber returns all wanted vehicles by number.
func (s *Store) VehiclesByNumber(number string) ([]Vehicle, error) {
	return s.db.VehiclesByNumber(number)
}

// WantedVehiclesByVIN returns all wanted vehicles by vin.
func (s *Store) VehiclesByVIN(vin string) ([]Vehicle, error) {
	return s.db.VehiclesByVIN(vin)
}

func (s *Store) RevisionByID(id string) (*Revision, error) {
	return s.db.RevisionByID(id)
}

// Vehicles returns list of vehicles ordered by theft_date with specified limit.
func (s *Store) Vehicles(limit int64) ([]Vehicle, error) {
	return s.db.Vehicles(limit)
}

func (s *Store) Revisions(limit int64) ([]Revision, error) {
	return s.db.Revisions(limit)
}
