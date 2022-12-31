package service

import (
	"context"

	"github.com/opencars/seedwork"

	"github.com/opencars/wanted/pkg/domain"
	"github.com/opencars/wanted/pkg/domain/model"
	"github.com/opencars/wanted/pkg/domain/query"
)

type CustomerService struct {
	vehicle  domain.VehicleRepository
	revision domain.RevisionRepository
	// producer schema.Producer
}

func NewCustomerService(v domain.VehicleRepository, r domain.RevisionRepository) *CustomerService {
	return &CustomerService{
		vehicle:  v,
		revision: r,
		// producer: producer,
	}
}

func (s *CustomerService) FindRevisionByID(ctx context.Context, id string) (*model.Revision, error) {
	return s.revision.FindByID(id)
}

func (s *CustomerService) ListRevisions(ctx context.Context, q *query.ListRevisions) ([]model.Revision, error) {
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	return s.revision.AllWithLimit(q.GetLimit())
}

func (s *CustomerService) ListByNumber(ctx context.Context, q *query.ListByNumber) ([]model.Vehicle, error) {
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	return s.vehicle.FindByNumber(q.Number)
}

func (s *CustomerService) ListByVIN(ctx context.Context, q *query.ListByVIN) ([]model.Vehicle, error) {
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	return s.vehicle.FindByVIN(q.VIN)
}

func (s *CustomerService) ListVehicles(ctx context.Context, q *query.ListVehicles) ([]model.Vehicle, error) {
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	return s.vehicle.AllWithLimit(q.GetLimit())
}
