package httpserver

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	e *echo.Echo
	// services *services.Services
}

func (s *EchoServer) Start(addr string) error            { return s.e.Start(addr) }
func (s *EchoServer) Shutdown(ctx context.Context) error { return s.e.Shutdown(ctx) }

func NewEchoServer() *EchoServer {
	e := echo.New()

	server := &EchoServer{
		e: e,
		// services: svcs,
	}

	server.setupRoutes()
	return server
}
func (s *EchoServer) setupRoutes() {
	apiV1 := s.e.Group("/api/v1")
	apiV1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})
}
