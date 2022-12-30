package service

import (
	"context"

	"github.com/opencars/seedwork"
	"github.com/opencars/wanted/pkg/domain"
	"github.com/opencars/wanted/pkg/domain/query"
)

type InternalService struct {
	r domain.VehicleRepository
}

func NewInternalService(r domain.VehicleRepository) *InternalService {
	return &InternalService{
		r: r,
	}
}

func (s *InternalService) Find(ctx context.Context, q *query.Find) (*query.FindResult, error) {
	if err := seedwork.ProcessQuery(q); err != nil {
		return nil, err
	}

	return s.r.Find(ctx, q)
}
