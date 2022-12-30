package http

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opencars/seedwork/httputil"

	"github.com/opencars/wanted/pkg/domain"
	"github.com/opencars/wanted/pkg/domain/query"
	"github.com/opencars/wanted/pkg/handler"
)

type VehicleAPI struct {
	svc domain.CustomerService
}

func (api *VehicleAPI) FindByVIN() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		q := query.ListByVIN{
			UserID:  httputil.UserIDFromContext(r.Context()),
			TokenID: httputil.TokenIDromContext(r.Context()),
			VIN:     mux.Vars(r)["vin"],
		}

		vehicles, err := api.svc.ListByVIN(r.Context(), &q)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(vehicles)
	}
}

func (api *VehicleAPI) FindByNumber() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		q := query.ListByNumber{
			UserID:  httputil.UserIDFromContext(r.Context()),
			TokenID: httputil.TokenIDromContext(r.Context()),
			Number:  mux.Vars(r)["number"],
		}

		vehicles, err := api.svc.ListByNumber(r.Context(), &q)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(vehicles)
	}
}

func (api *VehicleAPI) List() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		q := query.ListVehicles{
			UserID:  httputil.UserIDFromContext(r.Context()),
			TokenID: httputil.TokenIDromContext(r.Context()),
			Limit:   r.URL.Query().Get("limit"),
			Offset:  r.URL.Query().Get("offset"),
		}

		vehicles, err := api.svc.ListVehicles(r.Context(), &q)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(vehicles)
	}
}
