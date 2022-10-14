package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/opencars/seedwork/logger"
	"github.com/opencars/wanted/pkg/apiserver"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/store/postgres"
)

func main() {
	var configPath string

	flag.StringVar(&configPath, "config", "./config/config.yaml", "Path to the configuration file")

	flag.Parse()

	// Get configuration.
	conf, err := config.New(configPath)
	if err != nil {
		logger.Fatalf("failed to read config: %s", err)
	}

	store, err := postgres.New(conf)
	if err != nil {
		logger.Fatalf("failed to start postgres: %s", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		cancel()
	}()

	addr := ":8080"
	logger.Infof("Listening on %s...", addr)
	if err := apiserver.Start(ctx, addr, store); err != nil {
		logger.Fatalf("server failed: %s", err)
	}
}
