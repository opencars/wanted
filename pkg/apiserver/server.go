package apiserver

import (
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/store"
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

	core.Handle("/api/v1/wanted/revisions", s.Revision().All())
	core.Handle("/api/v1/wanted/revisions/stats", s.Revision().Stats())
	core.Handle("/api/v1/wanted/revisions/{id}", s.Revision().FindByID())

	core.Handle("/api/v1/wanted/swagger.yml", s.Swagger())

	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByNumber()).Queries("number", "{number}")
	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByVIN()).Queries("vin", "{vin}")
	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByRevisionID()).Queries("revision", "{revision}")
	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().All())
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})
	headers := handlers.AllowedHeaders([]string{"Api-Key", "X-Api-Key"})

	cors := handlers.CORS(origins, methods, headers)(s.router)
	compress := handlers.CompressHandler(cors)
	compress.ServeHTTP(w, r)
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

func (*server) Swagger() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./docs/swagger.yml")
	}
}

func (*server) limit(r *http.Request) (uint64, error) {
	limit, err := strconv.ParseUint(r.URL.Query().Get("limit"), 10, 64)
	if err != nil {
		return 100, nil
	}

	return limit, nil
}
