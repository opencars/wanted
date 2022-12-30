package main

import (
	"context"
	"flag"
	"os/signal"
	"strconv"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/opencars/seedwork/logger"

	"github.com/opencars/wanted/pkg/api/grpc"
	"github.com/opencars/wanted/pkg/config"
	"github.com/opencars/wanted/pkg/domain/service"
	"github.com/opencars/wanted/pkg/store/postgres"
)

func main() {
	cfg := flag.String("config", "config/config.yaml", "Path to the configuration file")
	port := flag.Int("port", 3000, "Port of the server")

	flag.Parse()

	conf, err := config.New(*cfg)
	if err != nil {
		logger.Fatalf("config: %v", err)
	}

	logger.NewLogger(logger.LogLevel(conf.Log.Level), conf.Log.Mode == "dev")

	store, err := postgres.New(conf)
	if err != nil {
		logger.Fatalf("store: %v", err)
	}

	svc := service.NewInternalService(store.Vehicle())

	addr := ":" + strconv.Itoa(*port)
	api := grpc.New(addr, svc)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof("Listening on %s...", nil)
	if err := api.Run(ctx); err != nil {
		logger.Fatalf("grpc: %v", err)
	}
}
