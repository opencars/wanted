package teststore

import (
	"sort"
	"time"

	"github.com/opencars/wanted/pkg/domain/model"
)

type RevisionRepository struct {
	store     *Store
	revisions map[string]*model.Revision
}

func (r *RevisionRepository) Create(revision *model.Revision) error {
	r.revisions[revision.ID] = revision

	return nil
}

func (r *RevisionRepository) FindByID(id string) (*model.Revision, error) {
	revision, ok := r.revisions[id]
	if !ok {
		return nil, model.ErrNotFound
	}

	return revision, nil
}

func (r *RevisionRepository) Last() (*model.Revision, error) {
	last := model.Revision{CreatedAt: time.Time{}}

	for _, v := range r.revisions {
		if v.CreatedAt.After(last.CreatedAt) {
			last = *v
		}
	}

	return &last, nil
}

func (r *RevisionRepository) All() ([]model.Revision, error) {
	revisions := make([]model.Revision, 0, len(r.revisions))

	for _, v := range r.revisions {
		revisions = append(revisions, *v)
	}

	return revisions, nil
}

func (r *RevisionRepository) AllWithLimit(_ uint64) ([]model.Revision, error) {
	revisions := make([]model.Revision, 0, len(r.revisions))

	for _, v := range r.revisions {
		revisions = append(revisions, *v)
	}

	return revisions, nil
}

func (r *RevisionRepository) AllIDs() ([]string, error) {
	revisions := make([]string, 0, len(r.revisions))

	for _, v := range r.revisions {
		revisions = append(revisions, v.ID)
	}

	return revisions, nil
}

// Stats returns aggregated revisions by month.
func (r *RevisionRepository) Stats() ([]model.RevisionStatMonth, error) {
	stats := make([]model.RevisionStatMonth, 0)

	revisions := make([]*model.Revision, 0, len(r.revisions))
	for _, v := range r.revisions {
		revisions = append(revisions, v)
	}

	// sort by descending.
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].CreatedAt.After(revisions[j].CreatedAt)
	})

	var current time.Month = -1
	for _, v := range r.revisions {
		if v.CreatedAt.Month() == current {
			stats[len(stats)-1].Added += v.Added
			stats[len(stats)-1].Removed += v.Removed
		} else {
			stats = append(stats, model.RevisionStatMonth{
				Month:   v.CreatedAt.Month(),
				Year:    v.CreatedAt.Year(),
				Added:   v.Added,
				Removed: v.Removed,
			})
			current = v.CreatedAt.Month()
		}
	}

	if len(stats) > 12 {
		stats = stats[:12]
	}

	return stats, nil
}
