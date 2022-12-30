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

	router.Handle("/wanted/revisions", s.Revision().All()).Methods("GET")
	router.Handle("/wanted/revisions/{id}", s.Revision().FindByID()).Methods("GET")

	router.Handle("/wanted/vehicles", s.Vehicle().FindByNumber()).Queries("number", "{number}").Methods("GET")
	router.Handle("/wanted/vehicles", s.Vehicle().FindByVIN()).Queries("vin", "{vin}").Methods("GET")
	router.Handle("/wanted/vehicles", s.Vehicle().List()).Methods("GET")
}

func (s *server) Revision() *RevisionAPI {
	return s.revisionAPI
}

func (s *server) Vehicle() *VehicleAPI {
	return s.vehicleAPI
}
