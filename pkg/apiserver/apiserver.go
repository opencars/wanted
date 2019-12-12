package apiserver

import (
	"net/http"

	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/store/postgres"
)

// Start starts the server with postgres store.
func Start(settings *config.Settings) error {
	store, err := postgres.New(settings)
	if err != nil {
		return err
	}

	srv := newServer(store)

	return http.ListenAndServe(":8080", srv)
}
