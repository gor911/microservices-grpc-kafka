package grpcserver

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/gor911/microservices-grpc-kafka/inventory-service/internal/controller/grpchandler"
	invpb "github.com/gor911/microservices-grpc-kafka/inventory-service/pkg/grpc/inventorypb"
	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

type Server struct {
	handler *grpchandler.Handlers
	grpc    *grpc.Server
	logger  *slog.Logger
	config  Config
}

func New(i *grpchandler.Handlers, logger *slog.Logger, c Config) *Server {
	return &Server{handler: i, logger: logger, config: c}
}

func (s *Server) Listen() error {
	lis, err := net.Listen("tcp", s.config.Addr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.grpc = grpc.NewServer()
	invpb.RegisterInventoryServer(s.grpc, s.handler)
	s.logger.Debug("grpc server listening on " + s.config.Addr)

	return s.grpc.Serve(lis)
}

func (s *Server) Shutdown() {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
