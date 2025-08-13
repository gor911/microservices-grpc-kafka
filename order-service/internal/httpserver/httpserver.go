package httpserver

import (
	"time"

	"github.com/valyala/fasthttp"
)

type Config struct {
	Addr string
}

type Server struct {
	server *fasthttp.Server
}

func New(handler fasthttp.RequestHandler, config Config) *Server {
	return &Server{
		server: &fasthttp.Server{
			Handler:      handler,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
		},
	}
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown()
}

func (s *Server) Listen(addr string) error {
	return s.server.ListenAndServe(addr)
}
