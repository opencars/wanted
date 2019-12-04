package main

import (
	"flag"
	"log"
	"strings"

	_ "github.com/lib/pq"

	"github.com/opencars/govdata"
	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/storage"
	"github.com/opencars/wanted/pkg/storage/postgres"
	"github.com/opencars/wanted/pkg/worker"
)

// ResourceID is a unique id of the resource in the government data provider.
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

	last, err := store.LastRevision()
	if err != nil {
		log.Fatal(err)
	}

	revisions := govdata.Subscribe(ResourceID, last.CreatedAt)
	for revision := range revisions {
		parts := strings.Split(revision.URL, "/")

		var vehicles []storage.WantedVehicle
		record := storage.Revision{
			ID:          parts[len(parts)-1],
			URL:         revision.URL,
			FileHashSum: revision.FileHashSum,
			CreatedAt:   revision.ResourceCreated.Time,
		}
		log.Println("Revision:", record.ID)

		body, err := govdata.ResourceRevision(resource.PackageID, ResourceID, record.ID)
		if err != nil {
			log.Fatal(revision.ID, err)
		}

		vehicles, record.Added, record.Removed, err = w.Parse(record.ID, bom.NewReader(body))
		if err == worker.ErrEmptyArr {
			continue
		}

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
