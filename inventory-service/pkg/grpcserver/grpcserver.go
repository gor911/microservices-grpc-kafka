package grpcserver

import (
	"errors"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type Config struct {
	Addr string
}

// RegisterFn attaches services to the grpc.Server.
// Example (in main):
//
//	s := grpcserver.New(":50051", func(gs *grpc.Server) {
//	    invpb.RegisterInventoryServer(gs, handler)
//	})
type RegisterFn func(*grpc.Server)

type Server struct {
	grpc   *grpc.Server
	logger *slog.Logger
	config Config
}

func New(register RegisterFn, logger *slog.Logger, c Config) (*Server, error) {
	if register == nil {
		return nil, errors.New("nil RegisterFn")
	}

	gs := grpc.NewServer()

	// register handlers now, before Start
	register(gs)

	return &Server{grpc: gs, logger: logger, config: c}, nil
}

func (s *Server) Listen() error {
	lis, err := net.Listen("tcp", s.config.Addr)

	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	s.logger.Debug("grpc server listening on " + s.config.Addr)

	return s.grpc.Serve(lis)
}

func (s *Server) Shutdown() {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
}
