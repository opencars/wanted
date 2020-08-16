package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/opencars/translit"

	"github.com/opencars/wanted/pkg/handler"
)

type VehicleAPI struct {
	server *server
}

func (api *VehicleAPI) FindByVIN() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		vin := mux.Vars(r)["vin"]
		transport, err := api.server.store.Vehicle().FindByVIN(vin)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&transport); err != nil {
			return err
		}

		return nil
	}
}

func (api *VehicleAPI) FindByNumber() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		number := mux.Vars(r)["number"]
		number = translit.ToUA(number)

		w.Header().Set("Content-Type", "application/json")

		transport, err := api.server.store.Vehicle().FindByNumber(number)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&transport); err != nil {
			return err
		}

		return nil
	}
}

func (api *VehicleAPI) FindByRevisionID() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["revision"]

		w.Header().Set("Content-Type", "application/json")

		transport, err := api.server.store.Vehicle().FindByRevisionID(id)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&transport); err != nil {
			return err
		}

		return nil
	}
}

func (api *VehicleAPI) All() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		limit, err := api.server.limit(r)
		if err != nil {
			return err
		}

		transport, err := api.server.store.Vehicle().AllWithLimit(limit)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&transport); err != nil {
			return err
		}

		return nil
	}
}
