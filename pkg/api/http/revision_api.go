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

type RevisionAPI struct {
	svc domain.CustomerService
}

// FindByID returns one revision found by it's unique id.
func (api *RevisionAPI) FindByID() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		id := mux.Vars(r)["id"]

		w.Header().Set("Content-Type", "application/json")

		revision, err := api.svc.FindRevisionByID(r.Context(), id)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(revision)
	}
}

// All returns last N revision from the store.
func (api *RevisionAPI) All() handler.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		w.Header().Set("Content-Type", "application/json")

		q := query.ListRevisions{
			UserID:  httputil.UserIDFromContext(r.Context()),
			TokenID: httputil.TokenIDromContext(r.Context()),
			Limit:   r.URL.Query().Get("limit"),
			Offset:  r.URL.Query().Get("offset"),
		}

		revisions, err := api.svc.ListRevisions(r.Context(), &q)
		if err != nil {
			return err
		}

		return json.NewEncoder(w).Encode(revisions)
	}
}

// Stats returns application parsing statistics for last 12 month.
// func (api *RevisionAPI) Stats() handler.Handler {
// 	return func(w http.ResponseWriter, r *http.Request) error {
// 		// w.Header().Set("Content-Type", "application/json")
// 		// stats, err := api.server.store.Revision().Stats()
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		// if err := json.NewEncoder(w).Encode(&stats); err != nil {
// 		// 	return err
// 		// }

// 		// return nil
// 	}
// }
