package apiserver

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/handler"
)

type RevisionAPI struct {
	server *server
}

// FindByID returns one revision found by it's unique id.
func (api *RevisionAPI) FinByID() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]

		w.Header().Set("Content-Type", "application/json")

		transport, err := api.server.store.Revision().FindByID(id)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&transport); err != nil {
			return err
		}

		return nil
	}
}

// All returns last N revision from the store.
func (api *RevisionAPI) All() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		limit, err := api.server.limit(r)
		if err != nil {
			return err
		}

		revisions, err := api.server.store.Revision().AllWithLimit(limit)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&revisions); err != nil {
			return err
		}

		return nil
	}
}

// Stats returns application parsing statistics for last 12 month.
func (api *RevisionAPI) Stats() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")
		stats, err := api.server.store.Revision().Stats()
		if err != nil {
			return err
		}

		if err := json.NewEncoder(w).Encode(&stats); err != nil {
			return err
		}

		return nil
	}
}
