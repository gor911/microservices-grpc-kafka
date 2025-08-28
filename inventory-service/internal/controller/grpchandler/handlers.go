package grpchandler

import (
	"context"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/dto"
	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/service"
	invpb "github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/grpc/inventorypb"
)

type Handlers struct {
	invpb.UnimplementedInventoryServer
	inventoryService *service.Inventory
}

func New(s *service.Inventory) *Handlers {
	return &Handlers{inventoryService: s}
}

func (h *Handlers) GetProducts(ctx context.Context, r *invpb.GetProductsRequest) (*invpb.GetProductsResponse, error) {
	var input dto.GetProductsInput

	input.ProductIds = make([]uint, len(r.GetProductIds()))

	for i, id := range r.GetProductIds() {
		input.ProductIds[i] = uint(id)
	}

	output, err := h.inventoryService.GetProducts(ctx, input)

	if err != nil {
		return nil, err
	}

	response := &invpb.GetProductsResponse{}

	for _, item := range output.Products {
		response.Products = append(response.Products, &invpb.Product{
			Id:    uint64(item.ID),
			Name:  item.Name,
			Price: item.Price,
		})
	}

	return response, nil
}
