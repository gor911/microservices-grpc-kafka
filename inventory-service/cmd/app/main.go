package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/config"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/controller/grpchandler"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/controller/kafkaconsumer"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/service"
	invpb "github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/grpc/inventorypb"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/grpcserver"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/logger"
	"google.golang.org/grpc"
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

	config.KafkaConsumer.Topic = "order.created"
	kafkaConsumer := kafkaconsumer.NewConsumer(config.KafkaConsumer, log, inventoryService.HandleProductStock)

	go func() {
		if err := kafkaConsumer.Consume(ctx); err != nil {
			panic(err)
		}
	}()

	gh := grpchandler.New(inventoryService)
	grpcServer, err := grpcserver.New(func(gs *grpc.Server) {
		invpb.RegisterInventoryServer(gs, gh)
	}, log, config.GRPC)

	if err != nil {
		return fmt.Errorf("could not create grpc server: %w", err)
	}

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

	if err := kafkaConsumer.Close(); err != nil {
		log.Error("could not close kafka consumer", "error", err)
	}

	grpcServer.Shutdown()

	log.Info("App is shutting down gracefully")

	return nil
}
