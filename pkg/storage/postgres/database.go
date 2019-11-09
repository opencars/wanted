package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/opencars/wanted/pkg/storage"
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

func (db *database) CreateOrUpdateVehicle(v *storage.WantedVehicle) error {
	_, err := db.db.NamedExec(`INSERT INTO vehicles (id, ovd, brand, kind, color, number, body_number, chassis_number, engine_number, status, theft_date, insert_date, revision_id) VALUES (:id, :ovd, :brand, :kind, :color, :number, :body_number, :chassis_number, :engine_number, :status, :theft_date, :insert_date, :revision_id) ON CONFLICT(id) DO UPDATE SET status = :status`, v)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}

func (db *database) CreateRevision(r *storage.Revision) error {
	_, err := db.db.NamedExec(`INSERT INTO revisions (id, name, url, file_hash_sum, removed, added, created_at) VALUES (:id, :name, :url, :file_hash_sum, :removed, :added, :created_at) ON CONFLICT(id) DO NOTHING`, r)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	return nil
}

func (db *database) AllVehicles() ([]storage.WantedVehicle, error) {
	vehicles := make([]storage.WantedVehicle, 0)
	if err := db.db.Select(&vehicles, "SELECT * FROM vehicles"); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}

func (db *database) AllRevisions() ([]storage.Revision, error) {
	revisions := make([]storage.Revision, 0)
	if err := db.db.Select(&revisions, "SELECT * FROM revisions"); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return revisions, nil
}

func (db *database) AllRevisionIDs() ([]string, error) {
	ids := make([]string, 0)
	if err := db.db.Select(&ids, "SELECT id FROM revisions"); err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return ids, nil
}

func (db *database) CreateOrUpdateVehicles(vv []storage.WantedVehicle) error {
	stmt, err := db.db.Preparex(`INSERT INTO vehicles (id, ovd, brand, kind, color, number, body_number, chassis_number, engine_number, status, theft_date, insert_date, revision_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT(id) DO UPDATE SET status = $10`)
	if err != nil {
		return err
	}

	for _, v := range vv {
		_, err := stmt.Exec(v.ID, v.OVD, v.Brand, v.Kind, v.Color, v.Number, v.BodyNumber, v.ChassisNumber, v.EngineNumber, v.Status, v.TheftDate, v.InsertDate, v.RevisionID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (db *database) CreateOrUpdateVehicles2(vv []storage.WantedVehicle) error {
	tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO vehicles (id, ovd, brand, kind, color, number, body_number, chassis_number, engine_number, status, theft_date, insert_date, revision_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT(id) DO UPDATE SET status = $10`)
	if err != nil {
		return err
	}

	for _, v := range vv {
		_, err := stmt.Exec(v.ID, v.OVD, v.Brand, v.Kind, v.Color, v.Number, v.BodyNumber, v.ChassisNumber, v.EngineNumber, v.Status, v.TheftDate, v.InsertDate, v.RevisionID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *database) CreateRevisionAndVehicles(revision *storage.Revision, vehicles []storage.WantedVehicle) error {
	tx, err := db.db.BeginTx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO revisions (id, name, url, file_hash_sum, removed, added, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT(id) DO NOTHING`,
		revision.ID, revision.Name, revision.URL,
		revision.FileHashSum, revision.Removed, revision.Added, revision.CreatedAt,
	)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO vehicles (id, ovd, brand, kind, color, number, body_number, chassis_number, engine_number, status, theft_date, insert_date, revision_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) ON CONFLICT(id) DO UPDATE SET status = $10`)
	if err != nil {
		return err
	}

	for _, v := range vehicles {
		_, err := stmt.Exec(v.ID, v.OVD, v.Brand, v.Kind, v.Color, v.Number, v.BodyNumber, v.ChassisNumber, v.EngineNumber, v.Status, v.TheftDate, v.InsertDate, v.RevisionID)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (db *database) VehiclesByNumber(number string) ([]storage.Vehicle, error) {
	vehicles := make([]storage.Vehicle, 0)

	err := db.db.Select(&vehicles, `SELECT * FROM vehicles WHERE number = $1`, number)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}

func (db *database) VehiclesByVIN(vin string) ([]storage.Vehicle, error) {
	vehicles := make([]storage.Vehicle, 0)

	err := db.db.Select(&vehicles, `SELECT * FROM vehicles WHERE body_number = $1 OR chassis_number = $1 OR engine_number = $1`, vin)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}

func (db *database) RevisionByID(id string) (*storage.Revision, error) {
	var revision storage.Revision

	err := db.db.Get(&revision, `SELECT * FROM revisions WHERE id LIKE $1`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return &revision, nil
}

// Vehicles returns list of vehicles ordered by theft_date with specified limit.
func (db *database) Vehicles(limit int64) ([]storage.Vehicle, error) {
	vehicles := make([]storage.Vehicle, 0)

	err := db.db.Select(&vehicles, `SELECT * FROM vehicles ORDER BY theft_date DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return vehicles, nil
}

// Revisions returns list of revisions ordered by created_at.
func (db *database) Revisions(limit int64) ([]storage.Revision, error) {
	revisions := make([]storage.Revision, 0)

	err := db.db.Select(&revisions, `SELECT * FROM revisions ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to select: %w", err)
	}

	return revisions, nil
}
