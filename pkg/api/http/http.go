package http

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/opencars/wanted/pkg/domain"
)

// Start starts the server with specified store.
func Start(ctx context.Context, addr string, svc domain.CustomerService) error {
	s := newServer(svc)

	srv := http.Server{
		Addr:    addr,
		Handler: handlers.LoggingHandler(os.Stdout, handlers.ProxyHeaders(s.router)),
	}

	errs := make(chan error)
	go func() {
		errs <- srv.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		ctxShutDown, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer func() {
			cancel()
		}()

		err := srv.Shutdown(ctxShutDown)
		if err != nil && err != http.ErrServerClosed {
			return err
		}

		return nil
	}
}
