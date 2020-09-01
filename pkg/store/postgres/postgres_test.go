package postgres_test

import (
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/opencars/wanted/pkg/config"
)

var settings *config.Settings

func TestMain(m *testing.M) {
	settings = &config.Settings{
		DB: config.Database{
			Host:     os.Getenv("DATABASE_HOST"),
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Name:     "wanted",
			SSLMode:  "disable",
		},
		Worker: config.Worker{
			ResourceID: "06e65b06-3120-4713-8003-7905a83f95f5",
		},
		Cleansing: config.Cleansing{
			Brand: config.BrandCleansing{
				Matchers: []config.Matcher{
					{
						Pattern: `^([BВ]{1}[АA]{1}[3З]{1})[\s-]*(.*)$`,
						Maker:   "ВАЗ",
						Model:   "$2",
					},
				},
			},
		},
	}

	if settings.DB.Host == "" {
		settings.DB.Host = "127.0.0.1"
	}

	os.Exit(m.Run())
}
