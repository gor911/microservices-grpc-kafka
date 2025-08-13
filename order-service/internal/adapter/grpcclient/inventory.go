package grpcclient

import (
	"context"
	"fmt"
	"time"

	invpb "github.com/gor911/microservices-grpc-kafka/order-service/pkg/grpc/inventorypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Config struct {
	Addr string
}

type InventoryClient struct {
	cc     *grpc.ClientConn
	client invpb.InventoryClient
}

func NewInventoryClient(c Config) (*InventoryClient, error) {
	cc, err := grpc.NewClient(
		c.Addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024)),
	)

	if err != nil {
		return nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	return &InventoryClient{cc: cc, client: invpb.NewInventoryClient(cc)}, nil
}

func (c *InventoryClient) Close() error {
	return c.cc.Close()
}

func (c *InventoryClient) GetProducts(ctx context.Context, ids []uint64) ([]*invpb.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	resp, err := c.client.GetProducts(ctx, &invpb.GetProductsRequest{ProductIds: ids})

	if err != nil {
		return nil, fmt.Errorf("client.GetProducts: %w", err)
	}

	return resp.GetProducts(), nil
}
