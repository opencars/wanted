package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/opencars/translit"
	"github.com/opencars/wanted/pkg/handler"
	"github.com/opencars/wanted/pkg/storage"
	"github.com/opencars/wanted/pkg/version"
)

type Server struct {
	store *storage.Store
}

func New(store *storage.Store) *Server {
	return &Server{
		store: store,
	}
}

func (_ *Server) Swagger(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./docs/swagger.yml")
}

func (_ *Server) Version(w http.ResponseWriter, r *http.Request) error {
	v := struct {
		Version string `json:"version"`
		Go      string `json:"go"`
	}{
		Version: version.Version,
		Go:      runtime.Version(),
	}

	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}

	return nil
}

func (srv *Server) VehiclesByVIN(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	vin := mux.Vars(r)["vin"]
	transport, err := srv.store.VehiclesByVIN(vin)
	if err != nil {
		return handler.NewError(http.StatusNotFound, "Record not found")
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		return err
	}

	return nil
}

func (srv *Server) VehiclesByNumber(w http.ResponseWriter, r *http.Request) error {
	number := mux.Vars(r)["number"]
	number = translit.ToUA(number)

	w.Header().Set("Content-Type", "application/json")

	transport, err := srv.store.VehiclesByNumber(number)
	if err != nil {
		return handler.NewError(http.StatusNotFound, "Record not found")
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		return err
	}

	return nil
}

func (srv *Server) VehiclesByRevisionID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["revision"]

	w.Header().Set("Content-Type", "application/json")

	transport, err := srv.store.VehiclesByRevisionID(id)
	if err != nil {
		return handler.NewError(http.StatusNotFound, "Record not found")
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		return err
	}

	return nil
}

func (srv *Server) Vehicles(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	transport, err := srv.store.Vehicles(100)
	if err != nil {
		return handler.NewError(http.StatusInternalServerError, "Unhealthy")
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		return err
	}

	return nil
}

func (srv *Server) RevisionByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	w.Header().Set("Content-Type", "application/json")

	transport, err := srv.store.RevisionByID(id)
	if err != nil {
		return handler.NewError(http.StatusNotFound, "Record not found")
	}

	if err := json.NewEncoder(w).Encode(&transport); err != nil {
		return err
	}

	return nil
}

func (srv *Server) Revisions(w http.ResponseWriter, r *http.Request) error {
	vars := r.URL.Query()

	w.Header().Set("Content-Type", "application/json")

	fmt.Println(vars.Get("limit"))
	var amount int64 = 100
	if vars.Get("limit") != "" {
		limit, err := strconv.ParseInt(vars.Get("limit"), 10, 64)
		if err != nil || limit <= 0 || limit > 100 {
			return handler.NewError(http.StatusBadRequest, "Limit is not valid")
		}

		// fmt.Println(limit)
		amount = limit
	}

	revisions, err := srv.store.Revisions(amount)
	if err != nil {
		return handler.NewError(http.StatusInternalServerError, "Unhealthy")
	}

	if err := json.NewEncoder(w).Encode(&revisions); err != nil {
		return err
	}

	return nil
}
