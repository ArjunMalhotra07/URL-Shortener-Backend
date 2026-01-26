package httpserver

import (
	"context"
	"net/http"
)

// StdServer is a thin wrapper around net/http.Server implementing Server.
type StdServer struct {
	server *http.Server
}

func NewStdServer(addr string, handler http.Handler) *StdServer {
	return &StdServer{server: &http.Server{
		Addr:    addr,
		Handler: handler,
	}}
}

func (s *StdServer) Start(addr string) error {
	if addr != "" {
		s.server.Addr = addr
	}
	return s.server.ListenAndServe()
}

func (s *StdServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
