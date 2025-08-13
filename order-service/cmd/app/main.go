package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gor911/microservices-grpc-kafka/order-service/config"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/grpcclient"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/kafkaproducer"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/controller/http"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/controller/publisher"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/httpserver"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/logger"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/service"
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

	inventoryClient, err := grpcclient.NewInventoryClient(config.GRPCClient)

	if err != nil {
		return fmt.Errorf("could not create inventory client: %w", err)
	}

	orderService := service.NewOrder(pg, log, inventoryClient)
	kafkaProducer := kafkaproducer.NewProducer(log, config.KafkaProducer)
	outboxService := service.NewOutbox(pg, log, kafkaProducer)
	kafkaPublisher := publisher.New(outboxService)

	go func() {
		if err := kafkaPublisher.Run(ctx); err != nil {
			panic(err)
		}
	}()

	httpServer := httpserver.New(http.BuildHandler(http.NewHandlers(orderService)), config.HTTP)

	go func() {
		if err := httpServer.Listen(config.HTTP.Addr); err != nil {
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

	if err := httpServer.Shutdown(); err != nil {
		log.Error("could not shutdown http server", "error", err)
	}

	if err := inventoryClient.Close(); err != nil {
		log.Error("could not close inventory client", "error", err)
	}

	if err := kafkaProducer.Close(); err != nil {
		log.Error("could not close kafka producer", "error", err)
	}

	if err := kafkaPublisher.Close(ctx); err != nil {
		log.Error("could not close kafka publisher worker", "error", err)
	}

	log.Info("App is shutting down gracefully")

	return nil
}
