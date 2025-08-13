package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/grpcclient"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/adapter/postgres"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/domain"
	"github.com/gor911/microservices-grpc-kafka/order-service/internal/dto"
	invpb "github.com/gor911/microservices-grpc-kafka/order-service/pkg/grpc/inventorypb"
	"github.com/valyala/fasthttp"
)

const topicOrderCreated = "order.created"

type Postgres interface {
	CreateOrder(ctx context.Context, order *domain.Order, outbox *domain.Outbox) error
	GetOutboxes(ctx context.Context) ([]*domain.Outbox, error)
	MarkOutboxesAsSent(ctx context.Context, ids []uint) error
}

type GRPCClient interface {
	GetProducts(ctx context.Context, ids []uint64) ([]*invpb.Product, error)
}

type Order struct {
	postgres   Postgres
	logger     *slog.Logger
	grpcClient GRPCClient
}

func NewOrder(postgres *postgres.Postgres, logger *slog.Logger, client *grpcclient.InventoryClient) *Order {
	return &Order{postgres: postgres, logger: logger, grpcClient: client}
}

func (o *Order) CreateOrder(ctx *fasthttp.RequestCtx, input dto.CreateOrderInput) (dto.CreateOrderOutput, error) {
	var output dto.CreateOrderOutput

	if input.UserID == 0 {
		return output, errors.New("invalid user id")
	}

	if len(input.Items) == 0 {
		return output, errors.New("invalid item count")
	}

	for _, item := range input.Items {
		if item.Qty == 0 {
			return output, errors.New("invalid quantity")
		}
	}

	productIds := make([]uint64, 0, len(input.Items))

	for _, item := range input.Items {
		productIds = append(productIds, uint64(item.ProductID))
	}

	products, err := o.grpcClient.GetProducts(ctx, productIds)

	if err != nil {
		return output, err
	}

	orderModel := domain.NewOrder(input, products)
	outboxModel := domain.NewOutbox(topicOrderCreated, input)

	if err := o.postgres.CreateOrder(ctx, orderModel, outboxModel); err != nil {
		return output, err
	}

	output.OrderID = orderModel.ID

	return output, nil
}
