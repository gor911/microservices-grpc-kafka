package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/config"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/controller/grpc"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/logger"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/service"
)

func main() {
	c, err := config.InitConfig()

	if err != nil {
		panic(err)
	}

	err = run(context.Background(), c)

	if err != nil {
		panic(err)
	}
}

func run(ctx context.Context, config config.Config) error {
	pg, err := postgres.New(ctx, config.Postgres)

	if err != nil {
		return fmt.Errorf("could not connect to postgres: %w", err)
	}

	log := logger.New(ctx, config.Logger)

	inventoryService := service.NewInventory(pg, log)

	grpcServer := grpc.New(inventoryService, log, config.GRPC)

	go func() {
		if err := grpcServer.Listen(); err != nil {
			panic(err)
		}
	}()

	log.Debug("App is running")

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	// closing all resources

	if err := pg.Close(ctx); err != nil {
		log.Error("could not close postgres connection", "error", err)
	}

	grpcServer.Shutdown()

	log.Info("App is shutting down gracefully")

	return nil
}
