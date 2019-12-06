package main

import (
	"flag"
	"log"

	"github.com/opencars/wanted/pkg/model"

	_ "github.com/lib/pq"

	"github.com/opencars/govdata"
	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/store/postgres"
	"github.com/opencars/wanted/pkg/worker"
)

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

	w := worker.New()

	resource, err := govdata.ResourceShow(conf.Worker.ResourceID)
	if err != nil {
		log.Fatal(err)
	}

	if err := w.Load(db); err != nil {
		log.Fatal(err)
	}

	last, err := db.Revision().Last()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Last revision:", last.ID)
	revisions := govdata.Subscribe(conf.Worker.ResourceID, last.CreatedAt)

	// Listen for new revisions.
	for revision := range revisions {
		record := model.RevisionFromGov(&revision)

		log.Println("Revision:", record.ID)

		body, err := govdata.ResourceRevision(resource.PackageID, conf.Worker.ResourceID, record.ID)
		if err != nil {
			log.Fatal(revision.ID, err)
		}

		reader, err := bom.NewReader(body)
		if err != nil {
			log.Fatal(err)
		}

		var vehicles []model.Vehicle
		vehicles, record.Added, record.Removed, err = w.Parse(record.ID, reader)
		if err == worker.ErrEmptyArr {
			continue
		}

		if err != nil {
			log.Fatalf("error %s duriding parsing of %s record", err, revision.ID)
		}

		log.Println("Added:", record.Added)
		log.Println("Removed:", record.Removed)

		if err := body.Close(); err != nil {
			log.Fatal(revision.ID, err)
		}

		// Save vehicles and revision.
		if err := db.Vehicle().Create(record, vehicles...); err != nil {
			log.Fatal(revision.ID, err)
		}
	}
}
