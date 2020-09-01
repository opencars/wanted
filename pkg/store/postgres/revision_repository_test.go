package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store"
	"github.com/opencars/wanted/pkg/store/postgres"
)

func TestRevisionRepository_Create(t *testing.T) {
	revision := model.TestRevision(t)

	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	assert.NoError(t, s.Revision().Create(revision))
	assert.NotNil(t, revision)
}

func TestRevisionRepository_FindByID(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	id := "example"
	_, err := s.Revision().FindByID(id)
	assert.Equal(t, store.ErrRecordNotFound, err)

	revision := model.TestRevision(t)
	revision.ID = id
	assert.NoError(t, s.Revision().Create(revision))

	found, err := s.Revision().FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, *revision, *found)
}

func TestRevisionRepository_Last(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	last, err := s.Revision().Last()
	assert.NoError(t, err)
	assert.Equal(t, time.Time{}, last.CreatedAt)

	revision := model.TestRevision(t)
	assert.NoError(t, s.Revision().Create(revision))

	last, err = s.Revision().Last()
	assert.NoError(t, err)
	assert.EqualValues(t, *revision, *last)
}

func TestRevisionRepository_All(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	revisions, err := s.Revision().All()
	assert.NoError(t, err)
	assert.Empty(t, revisions)

	revision := model.TestRevision(t)
	assert.NoError(t, s.Revision().Create(revision))

	revisions, err = s.Revision().All()
	assert.NoError(t, err)
	assert.NotEmpty(t, revisions)
	assert.Equal(t, *revision, revisions[0])
}

func TestRevisionRepository_AllWithLimit(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	revisions, err := s.Revision().All()
	assert.NoError(t, err)
	assert.Empty(t, revisions)

	revision := model.TestRevision(t)
	assert.NoError(t, s.Revision().Create(revision))

	revisions, err = s.Revision().AllWithLimit(100)
	assert.NoError(t, err)
	assert.NotEmpty(t, revisions)
	assert.Equal(t, *revision, revisions[0])
}

func TestRevisionRepository_AllIDs(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	ids, err := s.Revision().AllIDs()
	assert.NoError(t, err)
	assert.Empty(t, ids)

	revision := model.TestRevision(t)
	assert.NoError(t, s.Revision().Create(revision))

	ids, err = s.Revision().AllIDs()
	assert.NoError(t, err)
	assert.NotEmpty(t, ids)
	assert.Equal(t, revision.ID, ids[0])
}

func TestRevisionRepository_Stats(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("revisions")

	stats, err := s.Revision().Stats()
	assert.NoError(t, err)
	assert.Empty(t, stats)

	revision1 := model.TestRevision(t)
	revision2 := model.TestRevision(t)
	revision2.ID = "example_2"

	assert.NoError(t, s.Revision().Create(revision1))
	assert.NoError(t, s.Revision().Create(revision2))

	stats, err = s.Revision().Stats()
	assert.NoError(t, err)
	assert.NotEmpty(t, stats)
	assert.Equal(t, revision1.Added+revision2.Added, stats[0].Added)
	assert.Equal(t, revision1.Removed+revision2.Removed, stats[0].Removed)
	assert.Equal(t, revision1.CreatedAt.Month(), stats[0].Month)
}
