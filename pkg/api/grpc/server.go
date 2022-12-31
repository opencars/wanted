package grpc

import (
	"context"

	"github.com/opencars/grpc/pkg/wanted"
	"github.com/opencars/seedwork/logger"
	"github.com/opencars/wanted/pkg/domain/query"
)

type wantedHandler struct {
	wanted.UnimplementedServiceServer
	api *API
}

func (h *wantedHandler) Find(ctx context.Context, r *wanted.FindRequest) (*wanted.FindResponse, error) {
	q := query.Find{
		Numbers: r.Numbers,
		VINs:    r.Vins,
	}

	result, err := h.api.svc.Find(ctx, &q)
	if err != nil {
		return nil, handleErr(err)
	}

	logger.Debugf("vehicles: %+v", result)

	dtos := make([]*wanted.Vehicle, 0)
	for _, v := range result.Vehicles {
		dtos = append(dtos, fromDomain(&v))
	}

	logger.Debugf("dtos: %+v", dtos)

	return &wanted.FindResponse{
		Vehicles: dtos,
	}, nil
}
