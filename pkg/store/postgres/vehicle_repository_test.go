package postgres_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/store/postgres"
)

func TestVehicleRepository_Create(t *testing.T) {
	revision := model.TestRevision(t)
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle := model.TestVehicle(t)
	vehicle.RevisionID = revision.ID

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle}, nil))
	assert.NotNil(t, vehicle)
}

func TestVehicleRepository_FindByNumber(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle := model.TestVehicle(t)
	vehicles, err := s.Vehicle().FindByNumber(*vehicle.Number)
	assert.NoError(t, err)
	assert.NotNil(t, vehicles)
	assert.Len(t, vehicles, 0)

	revision := model.TestRevision(t)
	vehicle.RevisionID = revision.ID

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle}, nil))
	assert.NotNil(t, vehicle)

	vehicles, err = s.Vehicle().FindByNumber(*vehicle.Number)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, *vehicle, vehicles[0])
}

func TestVehicleRepository_FindByVIN(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle := model.TestVehicle(t)
	vehicles, err := s.Vehicle().FindByVIN(*vehicle.BodyNumber)
	assert.NoError(t, err)
	assert.NotNil(t, vehicles)
	assert.Len(t, vehicles, 0)

	revision := model.TestRevision(t)
	vehicle.RevisionID = revision.ID

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle}, nil))
	assert.NotNil(t, vehicle)

	vehicles, err = s.Vehicle().FindByVIN(*vehicle.BodyNumber)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, *vehicle, vehicles[0])
}

func TestVehicleRepository_FindByRevisionID(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle := model.TestVehicle(t)
	revision := model.TestRevision(t)
	vehicle.RevisionID = revision.ID

	vehicles, err := s.Vehicle().FindByRevisionID(vehicle.RevisionID)
	assert.NoError(t, err)
	assert.NotNil(t, vehicles)
	assert.Len(t, vehicles, 0)

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle}, nil))
	assert.NotNil(t, vehicle)

	vehicles, err = s.Vehicle().FindByRevisionID(vehicle.RevisionID)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, *vehicle, vehicles[0])
}

func TestVehicleRepository_AllWithLimit(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle1 := model.TestVehicle(t)
	vehicle2 := model.TestVehicle(t)
	vehicle2.CheckSum += "x"

	revision := model.TestRevision(t)
	vehicle1.RevisionID = revision.ID
	vehicle2.RevisionID = revision.ID

	vehicles, err := s.Vehicle().AllWithLimit(1)
	assert.NoError(t, err)
	assert.Empty(t, vehicles)

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle1, *vehicle2}, nil))

	vehicles, err = s.Vehicle().AllWithLimit(1)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, 1, len(vehicles))
	assert.Equal(t, *vehicle2, vehicles[0])

	vehicles, err = s.Vehicle().AllWithLimit(5)
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, 2, len(vehicles))
	assert.Equal(t, *vehicle2, vehicles[0])
	assert.Equal(t, *vehicle1, vehicles[1])
}

func TestVehicleRepository_All(t *testing.T) {
	s, teardown := postgres.TestDB(t, settings)
	defer teardown("vehicles", "revisions")

	vehicle1 := model.TestVehicle(t)
	vehicle2 := model.TestVehicle(t)
	vehicle2.CheckSum += "x"

	revision := model.TestRevision(t)
	vehicle1.RevisionID = revision.ID
	vehicle2.RevisionID = revision.ID

	vehicles, err := s.Vehicle().All()
	assert.NoError(t, err)
	assert.Empty(t, vehicles)

	assert.NoError(t, s.Vehicle().Create(revision, []model.Vehicle{*vehicle1, *vehicle2}, nil))

	vehicles, err = s.Vehicle().All()
	assert.NoError(t, err)
	assert.NotEmpty(t, vehicles)
	assert.Equal(t, 2, len(vehicles))
	assert.Equal(t, *vehicle1, vehicles[0])
	assert.Equal(t, *vehicle2, vehicles[1])
}
