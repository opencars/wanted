package apiserver

import (
	"encoding/json"
	"net/http"
	"runtime"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/handler"
	"github.com/opencars/wanted/pkg/store"
	"github.com/opencars/wanted/pkg/version"
)

type server struct {
	router *mux.Router
	store  store.Store

	revisionAPI *RevisionAPI
	vehicleAPI  *VehicleAPI
}

func newServer(store store.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	core := s.router.Methods("GET", "OPTIONS").Subrouter()

	core.Handle("/wanted/revisions", s.Revision().All())
	core.Handle("/wanted/revisions/stats", s.Revision().Stats())
	core.Handle("/wanted/revisions/{id}", s.Revision().FinByID())

	core.Handle("/wanted/swagger.yml", s.Swagger())
	core.Handle("/wanted/version", handler.Handler(s.Version))

	core.Handle("/wanted/vehicles", s.Vehicle().FindByNumber()).Queries("number", "{number}")
	core.Handle("/wanted/vehicles", s.Vehicle().FindByVIN()).Queries("vin", "{vin}")
	core.Handle("/wanted/vehicles", s.Vehicle().FindByRevisionID()).Queries("revision", "{revision}")
	core.Handle("/wanted/vehicles", s.Vehicle().All())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Api-Key"})

	cors := handlers.CORS(origins, methods, headers)(s.router)
	cors.ServeHTTP(w, r)
}

func (s *server) Revision() *RevisionAPI {
	if s.revisionAPI != nil {
		return s.revisionAPI
	}

	s.revisionAPI = &RevisionAPI{
		server: s,
	}

	return s.revisionAPI
}

func (s *server) Vehicle() *VehicleAPI {
	if s.vehicleAPI != nil {
		return s.vehicleAPI
	}

	s.vehicleAPI = &VehicleAPI{
		server: s,
	}

	return s.vehicleAPI
}

func (_ *server) Swagger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yml")
	}
}

func (_ *server) Version(w http.ResponseWriter, r *http.Request) error {
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
