package main

import (
	"flag"

	_ "github.com/lib/pq"
	"github.com/opencars/govdata"

	"github.com/opencars/wanted/pkg/bom"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/logger"
	"github.com/opencars/wanted/pkg/model"
	"github.com/opencars/wanted/pkg/store/postgres"
	"github.com/opencars/wanted/pkg/worker"
)

func main() {
	var path string

	flag.StringVar(&path, "config", "./config/config.yaml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(path)
	if err != nil {
		logger.Fatal(err)
	}

	// Register postgres adapter.
	db, err := postgres.New(conf)
	if err != nil {
		logger.Fatal(err)
	}

	resource, err := govdata.ResourceShow(conf.Worker.ResourceID)
	if err != nil {
		logger.Fatal(err)
	}

	w := worker.New()
	if err := w.Load(db); err != nil {
		logger.Fatal(err)
	}

	last, err := db.Revision().Last()
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Last revision: %s", last.ID)
	revisions := govdata.Subscribe(conf.Worker.ResourceID, last.CreatedAt)

	// Listen for new revisions.
	for revision := range revisions {
		record := model.RevisionFromGov(&revision)

		logger.WithFields(logger.Fields{
			"revision": record.ID,
		}).Info("Started processing revision")

		body, err := govdata.ResourceRevision(resource.PackageID, conf.Worker.ResourceID, record.ID)
		if err != nil {
			logger.WithFields(logger.Fields{
				"revision": record.ID,
			}).Fatal(err)
		}

		reader, err := bom.NewReader(body)
		if err != nil {
			logger.WithFields(logger.Fields{
				"revision": revision,
				"err":      err,
			}).Error("Broken bom encoding. Skipped")
			continue
		}

		added, removed, err := w.Parse(record.ID, reader)
		if err == worker.ErrEmptyArr {
			logger.WithFields(logger.Fields{
				"revision": revision,
			}).Error("Revision is empty. Skipped")
			continue
		}

		if err != nil {
			logger.WithFields(logger.Fields{
				"revision": revision,
				"err":      err,
			}).Error("Revision is broken. Skipped")
			continue
		}

		if err := body.Close(); err != nil {
			logger.WithFields(logger.Fields{
				"revision": revision,
				"err":      err,
			}).Error("Failed to close body. Skipped")
			continue
		}

		record.Added = len(added)
		record.Removed = len(removed)

		// Save vehicles and revision.
		if err := db.Vehicle().Create(record, added, removed); err != nil {
			logger.WithFields(logger.Fields{
				"revision": revision,
			}).Fatal(err)
		}

		logger.WithFields(logger.Fields{
			"revision": record.ID,
			"added":    len(added),
			"removed":  len(removed),
		}).Info("Finished processing revision")
	}
}
