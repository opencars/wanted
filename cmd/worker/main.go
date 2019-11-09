package main

import (
	"flag"
	"log"
	"strings"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/govdata"
	"github.com/opencars/wanted/pkg/storage"
	"github.com/opencars/wanted/pkg/storage/postgres"
	"github.com/opencars/wanted/pkg/worker"
)

const ResourceID = "06e65b06-3120-4713-8003-7905a83f95f5"

func main() {
	var path string

	flag.StringVar(&path, "config", "./config/config.toml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(path)
	if err != nil {
		log.Fatal(err)
	}

	// Register postgres adapter.
	db, err := postgres.New(conf.DB.Host, conf.DB.Port, conf.DB.User, conf.DB.Password, conf.DB.Name)
	if err != nil {
		log.Fatal(err)
	}

	store := storage.New(db)
	w := worker.New()

	resource, err := govdata.ResourceShow(ResourceID)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.Load(store); err != nil {
		log.Fatal(err)
	}

	ids, err := store.AllRevisionIDs()
	if err != nil {
		log.Fatal(err)
	}

	revisions := govdata.Subscribe(ResourceID, ids...)
	for revision := range revisions {
		parts := strings.Split(revision.URL, "/")

		var vehicles []storage.WantedVehicle
		record := storage.Revision{
			ID:          parts[len(parts)-1],
			Name:        revision.Name,
			URL:         revision.URL,
			FileHashSum: revision.FileHashSum,
			CreatedAt:   revision.ResourceCreated,
		}

		body, err := govdata.ResourceRevision(resource.PackageID, ResourceID, record.ID)
		if err != nil {
			log.Fatal(revision.ID, err)
		}

		log.Println("Revision:", record.ID)

		vehicles, record.Added, record.Removed, err = w.Parse(record.ID, bom.NewReader(body))
		if err != nil {
			log.Fatal(revision.ID, err)
		}

		log.Println("Added:", record.Added)
		log.Println("Removed:", record.Removed)
		log.Println("Changes:", len(vehicles))

		if err := body.Close(); err != nil {
			log.Fatal(revision.ID, err)
		}

		w.Fix(vehicles)

		if err := store.CreateRevisionAndVehicles(&record, vehicles); err != nil {
			log.Fatal(revision.ID, err)
		}
	}
}
