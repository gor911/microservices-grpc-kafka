package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/domain"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/dto"
	"github.com/segmentio/kafka-go"
)

type Postgres interface {
	GetProducts(ctx context.Context, ids []uint) ([]*domain.Product, error)
}

type Inventory struct {
	postgres Postgres
	logger   *slog.Logger
}

func NewInventory(postgres *postgres.Postgres, logger *slog.Logger) *Inventory {
	return &Inventory{postgres: postgres, logger: logger}
}

func (i *Inventory) GetProducts(ctx context.Context, input dto.GetProductsInput) (dto.GetProductsOutput, error) {
	var output dto.GetProductsOutput

	if len(input.ProductIds) == 0 {
		return output, errors.New("no products found")
	}

	products, err := i.postgres.GetProducts(ctx, input.ProductIds)

	if err != nil {
		return output, err
	}

	if len(products) != len(input.ProductIds) {
		return output, errors.New("product count mismatch")
	}

	for _, product := range products {
		output.Products = append(output.Products, dto.Product{ID: product.ID, Name: product.Name, Price: product.Price})
	}

	return output, nil
}

func (i *Inventory) HandleProductStock(ctx context.Context, msg kafka.Message) error {
	i.logger.Debug("Handling product stock")

	var orderCreated domain.OrderCreated

	if err := json.Unmarshal(msg.Value, &orderCreated); err != nil {
		return fmt.Errorf("invalid order created payload")
	}

	var productIds []uint

	for _, item := range orderCreated.Items {
		productIds = append(productIds, item.ProductID)
	}

	products, err := i.postgres.GetProducts(ctx, productIds)

	if err != nil {
		return err
	}

	if len(products) != len(productIds) {
		// todo produce order.failed
		return nil
	}

	// todo produce order.reserved

	return nil
}
