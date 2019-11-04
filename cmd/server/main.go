package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/server"
	"github.com/opencars/wanted/pkg/server/storage"
	"github.com/opencars/wanted/pkg/server/storage/postgres"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.toml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Register postgres adapter.
	db, err := postgres.New(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	store := storage.New(db)
	handler := server.New(store)

	router := mux.NewRouter()

	// router.HandleFunc("/version", handler.Version).Methods("GET")
	// router.HandleFunc("/health", handler.Health).Methods("GET")

	router.HandleFunc("/wanted-transport", handler.ByNumber).Queries("number", "{number}").Methods("GET")
	router.HandleFunc("/wanted-transport", handler.ByVIN).Queries("vin", "{vin}").Methods("GET")

	server := http.Server{
		Addr:    ":8080",
		Handler: handlers.LoggingHandler(os.Stdout, router),
	}

	log.Println("Listening on port 8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
