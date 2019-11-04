package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/server/storage"
)

type Transport struct {
	store *storage.Store
}

func New(store *storage.Store) *Transport {
	return &Transport{
		store: store,
	}
}

func (wt *Transport) ByVIN(w http.ResponseWriter, r *http.Request) {
	vin := mux.Vars(r)["vin"]

	w.Header().Set("Content-Type", "application/json")

	transport, err := wt.store.WantedVehiclesByVIN(vin)
	if err != nil {
		http.Error(w, "failed", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		http.Error(w, "failed", 500)
		return
	}
}

func (wt *Transport) ByNumber(w http.ResponseWriter, r *http.Request) {
	number := mux.Vars(r)["number"]

	w.Header().Set("Content-Type", "application/json")

	transport, err := wt.store.WantedVehiclesByNumber(number)
	if err != nil {
		http.Error(w, "failed", 500)
		return
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		http.Error(w, "failed", 500)
		return
	}
}
