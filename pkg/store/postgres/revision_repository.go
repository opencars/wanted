package postgres

import (
	"database/sql"
	"time"

	"github.com/opencars/wanted/pkg/store"

	"github.com/opencars/wanted/pkg/model"
)

type RevisionRepository struct {
	store *Store
}

func (r *RevisionRepository) Create(revision *model.Revision) error {
	_, err := r.store.db.NamedExec(`INSERT INTO revisions (id, url, file_hash_sum, removed, added, created_at) VALUES (:id, :url, :file_hash_sum, :removed, :added, :created_at) ON CONFLICT(id) DO NOTHING`, revision)
	if err != nil {
		return err
	}

	return nil
}

func (r *RevisionRepository) FindByID(id string) (*model.Revision, error) {
	var revision model.Revision

	err := r.store.db.Get(&revision, `SELECT * FROM revisions WHERE id LIKE $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	revision.CreatedAt = revision.CreatedAt.UTC()
	return &revision, nil
}

func (r *RevisionRepository) Last() (*model.Revision, error) {
	revisions := make([]model.Revision, 0)

	err := r.store.db.Select(&revisions, `SELECT * FROM revisions ORDER BY created_at DESC LIMIT 1`)
	if err == sql.ErrNoRows || len(revisions) == 0 {
		return &model.Revision{CreatedAt: time.Time{}}, nil
	}

	if err != nil {
		return nil, err
	}

	revisions[0].CreatedAt = revisions[0].CreatedAt.UTC()
	return &revisions[0], nil
}

func (r *RevisionRepository) All() ([]model.Revision, error) {
	revisions := make([]model.Revision, 0)
	if err := r.store.db.Select(&revisions, `SELECT * FROM revisions`); err != nil {
		return nil, err
	}

	for i := range revisions {
		revisions[i].CreatedAt = revisions[i].CreatedAt.UTC()
	}

	return revisions, nil
}

func (r *RevisionRepository) AllIDs() ([]string, error) {
	ids := make([]string, 0)
	if err := r.store.db.Select(&ids, `SELECT id FROM revisions`); err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *RevisionRepository) Stats() ([]model.RevisionStatMonth, error) {
	stats := make([]model.RevisionStatMonth, 0, 1000)

	query := `SELECT date_part('month', created_at) as month, date_part('year', created_at) as year, sum(removed) as removed, sum(added) as added
FROM revisions
GROUP BY year, month
ORDER BY year DESC, month DESC
LIMIT 12`

	if err := r.store.db.Select(&stats, query); err != nil {
		return nil, err
	}

	return stats, nil
}
