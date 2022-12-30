package postgres

import (
	"context"
	"database/sql"

	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/domain/query"
)

type VehicleRepository struct {
	store *Store
}

func (r *VehicleRepository) All() ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, kind, maker, model, color, number,
				body_number, chassis_number, engine_number,
				status, theft_date, insert_date, revision_id
		FROM vehicles ORDER BY theft_date`,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return vehicles, nil
}

func (r *VehicleRepository) CreateOrUpdateAll(vv []model.Vehicle) error {
	stmt, err := r.store.db.PrepareNamed(`
		INSERT INTO vehicles (
			id, ovd, brand, maker, model, kind, color,
			number, body_number, chassis_number, engine_number,
			status, theft_date, insert_date, revision_id
		) VALUES (
			:id, :ovd, :brand, :maker, :model, :kind, :color,
			:number, :body_number, :chassis_number, :engine_number,
			:status, :theft_date, :insert_date, :revision_id
		)
		ON CONFLICT(id) DO UPDATE SET status = :status`)
	if err != nil {
		return err
	}

	for _, v := range vv {
		v.BeforeCreate(r.store.clean)
		_, err := stmt.Exec(v)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *VehicleRepository) FindByNumber(number string) ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, maker, model, kind, color, number,
				body_number, chassis_number, engine_number,
				status, theft_date, insert_date, revision_id
		FROM vehicles
		WHERE number = $1`,
		number,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return vehicles, nil
}

func (r *VehicleRepository) FindByVIN(vin string) ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, maker, model, kind, color, number,
				body_number, chassis_number, engine_number,
				status, theft_date, insert_date, revision_id
		FROM vehicles
		WHERE body_number = $1 OR chassis_number = $1 OR engine_number = $1`,
		vin,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return vehicles, nil
}

func (r *VehicleRepository) FindByRevisionID(id string) ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, kind, maker, model, color, number,
				body_number, chassis_number, engine_number,
				status, theft_date, insert_date, revision_id
		FROM vehicles
		WHERE revision_id = $1`,
		id,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return vehicles, nil
}

func (r *VehicleRepository) AllWithLimit(limit uint64) ([]model.Vehicle, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, kind, maker, model, color, number,
					  body_number, chassis_number, engine_number,
					  status, theft_date, insert_date, revision_id
			FROM vehicles
			ORDER BY theft_date DESC, id DESC
			LIMIT $1`,
		limit,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return vehicles, nil
}

func (r *VehicleRepository) Create(revision *model.Revision, added []model.Vehicle, removed []string) error {
	tx, err := r.store.db.BeginTxx(context.Background(), &sql.TxOptions{})
	if err != nil {
		return err
	}
	_, err = tx.NamedExec(`
		INSERT INTO revisions (id, url, file_hash_sum, removed, added, created_at)
		VALUES (:id, :url, :file_hash_sum, :removed, :added, :created_at)
		ON CONFLICT(id) DO NOTHING`,
		revision,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	stmt, err := tx.PrepareNamed(`
		INSERT INTO vehicles (
			id, ovd, brand, maker, model, kind, color, number,
			body_number, chassis_number, engine_number, status,
			theft_date, insert_date, revision_id
		)
		VALUES (
			:id, :ovd, :brand, :maker, :model, :kind, :color, :number,
			:body_number, :chassis_number, :engine_number, :status,
			:theft_date, :insert_date, :revision_id
		)
        ON CONFLICT(id) DO
        UPDATE SET status = :status`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	for _, v := range added {
		v.BeforeCreate(r.store.clean)
		if _, err := stmt.Exec(v); err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	query := `UPDATE vehicles SET status = $2 WHERE id = $1`
	for _, id := range removed {
		_, err := tx.Exec(query, id, model.StatusRemoved)
		if err != nil {
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *VehicleRepository) Find(ctx context.Context, q *query.Find) (*query.FindResult, error) {
	vehicles := make([]model.Vehicle, 0)

	err := r.store.db.Select(&vehicles,
		`SELECT id, ovd, brand, maker, model, kind, color, number,
				body_number, chassis_number, engine_number,
				status, theft_date, insert_date, revision_id
		FROM vehicles
		WHERE body_number IN $1 OR chassis_number IN $1 OR engine_number IN $1 OR number IN $2`,
		q.VINs, q.Numbers,
	)
	if err != nil {
		return nil, err
	}

	for i := range vehicles {
		vehicles[i].InsertDate = vehicles[i].InsertDate.UTC()
	}

	return &query.FindResult{
		Vehicles: vehicles,
	}, nil
}
