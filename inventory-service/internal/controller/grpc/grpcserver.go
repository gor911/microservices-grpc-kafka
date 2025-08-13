package grpc

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/service"
	invpb "github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/grpc/inventorypb"
	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

type Server struct {
	invpb.UnimplementedInventoryServer
	inventoryService *service.Inventory
	grpc             *grpc.Server
	logger           *slog.Logger
	config           Config
}

func New(i *service.Inventory, logger *slog.Logger, c Config) *Server {
	return &Server{inventoryService: i, logger: logger, config: c}
}

func (s *Server) Listen() error {
	lis, err := net.Listen("tcp", s.config.Addr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpc = grpc.NewServer()
	invpb.RegisterInventoryServer(s.grpc, s)
	s.logger.Debug("grpc server listening on " + s.config.Addr)

	return s.grpc.Serve(lis)
}

func (s *Server) Shutdown() {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
