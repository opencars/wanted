package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/api"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/handler"
	"github.com/opencars/wanted/pkg/storage"
	"github.com/opencars/wanted/pkg/storage/postgres"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.toml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		log.Fatal(err)
	}

	// Register postgres adapter.
	db, err := postgres.New(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	api := api.New(store)

	router := mux.NewRouter().PathPrefix("/wanted/").Subrouter()

	core := router.Methods("GET", "OPTIONS").Subrouter()
	core.HandleFunc("/swagger.yml", api.Swagger)
	core.Handle("/version", handler.Handler(api.Version))

	core.Handle("/vehicles", handler.Handler(api.VehiclesByNumber)).Queries("number", "{number}")
	core.Handle("/vehicles", handler.Handler(api.VehiclesByVIN)).Queries("vin", "{vin}")
	core.Handle("/vehicles", handler.Handler(api.VehiclesByRevisionID)).Queries("revision", "{revision}")
	core.Handle("/vehicles", handler.Handler(api.Vehicles))

	core.Handle("/revisions", handler.Handler(api.Revisions))
	core.Handle("/revisions/{id}", handler.Handler(api.RevisionByID))

	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "OPTIONS"})

	cors := handlers.CORS(origins, methods)(router)
	srv := http.Server{
		Addr:    ":8080",
		Handler: handlers.LoggingHandler(os.Stdout, cors),
	}

	log.Println("Listening on port 8080...")
	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}
