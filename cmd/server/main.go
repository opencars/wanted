package main

import (
	"flag"

	_ "github.com/lib/pq"

	"github.com/opencars/wanted/pkg/apiserver"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/logger"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.toml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatal(err)
	}

	if err := apiserver.Start(conf); err != nil {
		logger.Fatal(err)
	}
}
