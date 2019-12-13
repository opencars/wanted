package teststore_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store/teststore"
)

func TestRevisionRepository_Create(t *testing.T) {
	s := teststore.New()

	revision := model.TestRevision(t)

	assert.NoError(t, s.Revision().Create(revision))
	assert.NotNil(t, revision)
}

func TestRevisionRepository_FindByID(t *testing.T) {
	s := teststore.New()

	id := "example"
	_, err := s.Revision().FindByID(id)
	assert.Error(t, err)

	revision := model.TestRevision(t)
	revision.ID = id
	assert.NoError(t, s.Revision().Create(revision))

	found, err := s.Revision().FindByID(id)
	assert.NoError(t, err)
	assert.Equal(t, revision, found)
}

func TestRevisionRepository_Last(t *testing.T) {
	s := teststore.New()

	last, err := s.Revision().Last()
	assert.NoError(t, err)
	assert.Equal(t, time.Time{}, last.CreatedAt)

	revision := model.TestRevision(t)
	assert.NoError(t, s.Revision().Create(revision))

	last, err = s.Revision().Last()
	assert.NoError(t, err)
	assert.Equal(t, revision, last)
}

func TestRevisionRepository_All(t *testing.T) {
	s := teststore.New()

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

func TestRevisionRepository_AllIDs(t *testing.T) {
	s := teststore.New()

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
	s := teststore.New()

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
