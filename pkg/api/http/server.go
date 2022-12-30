package http

import (
	"github.com/gorilla/mux"
	"github.com/opencars/seedwork/httputil"
	"github.com/opencars/wanted/pkg/domain"
)

type server struct {
	router *mux.Router

	revisionAPI *RevisionAPI
	vehicleAPI  *VehicleAPI
}

func newServer(svc domain.CustomerService) *server {
	s := &server{
		router: mux.NewRouter(),
		revisionAPI: &RevisionAPI{
			svc: svc,
		},
		vehicleAPI: &VehicleAPI{
			svc: svc,
		},
	}

	s.configureRouter()

	return s
}

func (s *server) configureRouter() {
	router := s.router.PathPrefix("/api/v1/").Subrouter()
	router.Use(
		httputil.CustomerTokenMiddleware(),
	)

	router.Handle("/api/v1/wanted/revisions", s.Revision().All())
	router.Handle("/api/v1/wanted/revisions/{id}", s.Revision().FindByID())

	router.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByNumber()).Queries("number", "{number}")
	router.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByVIN()).Queries("vin", "{vin}")
	router.Handle("/api/v1/wanted/vehicles", s.Vehicle().List())
}

func (s *server) Revision() *RevisionAPI {
	return s.revisionAPI
}

func (s *server) Vehicle() *VehicleAPI {
	return s.vehicleAPI
}
