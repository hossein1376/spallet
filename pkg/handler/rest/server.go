package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/hossein1376/spallet/pkg/application/service"
)

type Server struct {
	srv *http.Server
}

func NewServer(addr string, services *service.Services) *Server {
	mux := routes(services)

	return &Server{
		srv: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  30 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
