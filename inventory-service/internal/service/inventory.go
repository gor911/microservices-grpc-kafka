package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/domain"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/dto"
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
