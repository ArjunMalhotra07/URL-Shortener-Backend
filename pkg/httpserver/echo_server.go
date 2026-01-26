package httpserver

import (
	"context"
	"net/http"

	"url_shortner_backend/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoServer struct {
	e    *echo.Echo
	svcs *internal.AppServices
}

func (s *EchoServer) Start(addr string) error            { return s.e.Start(addr) }
func (s *EchoServer) Shutdown(ctx context.Context) error { return s.e.Shutdown(ctx) }

func NewEchoServer(svcs *internal.AppServices) *EchoServer {
	e := echo.New()

	// CORS middleware - allow credentials (cookies) from frontend
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowCredentials: true,
	}))

	server := &EchoServer{
		e:    e,
		svcs: svcs,
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

	// Short URL routes
	apiV1.POST("/shorten", s.svcs.ShortURL.CreateShortURL)
	apiV1.GET("/my-urls", s.svcs.ShortURL.GetMyURLs)
	apiV1.GET("/:code", s.svcs.ShortURL.GetOriginalURL)
}
