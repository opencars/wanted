package grpc

import (
	"context"

	"github.com/opencars/grpc/pkg/wanted"
	"github.com/opencars/wanted/pkg/domain/query"
)

type wantedHandler struct {
	wanted.UnimplementedServiceServer
	api *API
}

func (h *wantedHandler) Find(ctx context.Context, r *wanted.FindRequest) (*wanted.FindResponse, error) {
	q := query.Find{
		Numbers: r.Number,
		VINs:    r.Vin,
	}

	result, err := h.api.svc.Find(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	dtos := make([]*wanted.Vehicle, 0)
	for _, v := range result.Vehicles {
		dtos = append(dtos, fromDomain(&v))
	}

	return &wanted.FindResponse{
		Vehicles: dtos,
	}, nil
}
