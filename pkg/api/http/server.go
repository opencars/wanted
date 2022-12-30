package http

import (
	"github.com/gorilla/mux"
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
	core := s.router.Methods("GET", "OPTIONS").Subrouter()

	core.Handle("/api/v1/wanted/revisions", s.Revision().All())
	core.Handle("/api/v1/wanted/revisions/{id}", s.Revision().FindByID())

	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByNumber()).Queries("number", "{number}")
	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().FindByVIN()).Queries("vin", "{vin}")
	core.Handle("/api/v1/wanted/vehicles", s.Vehicle().List())
}

func (s *server) Revision() *RevisionAPI {
	return s.revisionAPI
}

func (s *server) Vehicle() *VehicleAPI {
	return s.vehicleAPI
}
