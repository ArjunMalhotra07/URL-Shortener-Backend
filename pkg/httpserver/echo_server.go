package httpserver

import (
	"context"
	"net/http"
	"time"

	"url_shortner_backend/internal"
	authmw "url_shortner_backend/pkg/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"

	"url_shortner_backend/pkg/jwt"
)

type EchoServer struct {
	e    *echo.Echo
	svcs *internal.AppServices
	jwt  *jwt.JWTManager
}

func (s *EchoServer) Start(addr string) error            { return s.e.Start(addr) }
func (s *EchoServer) Shutdown(ctx context.Context) error { return s.e.Shutdown(ctx) }

func NewEchoServer(svcs *internal.AppServices, jwtMgr *jwt.JWTManager) *EchoServer {
	e := echo.New()

	// CORS middleware - allow credentials (cookies) from frontend
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // Required because you use credentials: "include"
		AllowOrigins: []string{
			"http://localhost:5173",
			"http://192.168.1.3:5173",
			"https://tinyclk.com",
			"http://tinyclk.com",
		},
		AllowMethods:  []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		ExposeHeaders: []string{"Set-Cookie"},
	}))

	// Auth middleware - extracts user from JWT (does not block)
	e.Use(authmw.AuthMiddleware(jwtMgr))

	server := &EchoServer{
		e:    e,
		svcs: svcs,
		jwt:  jwtMgr,
	}

	server.setupRoutes()
	return server
}

func (s *EchoServer) setupRoutes() {
	// Rate limiters
	// General API: 20 requests/second with burst of 40
	apiRateLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(20),
			Burst:     40,
			ExpiresIn: 3 * time.Minute,
		},
	))

	// Auth routes: stricter limit - 5 requests/second with burst of 10 (prevent brute force)
	authRateLimiter := middleware.RateLimiter(middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate:      rate.Limit(5),
			Burst:     10,
			ExpiresIn: 3 * time.Minute,
		},
	))

	apiV1 := s.e.Group("/api/v1", apiRateLimiter)
	apiV1.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "healthy",
		})
	})

	// Auth routes (with stricter rate limit)
	auth := apiV1.Group("/auth", authRateLimiter)
	auth.POST("/signup", s.svcs.Auth.Signup)
	auth.POST("/login", s.svcs.Auth.Login)
	auth.POST("/refresh", s.svcs.Auth.Refresh)
	auth.POST("/logout", s.svcs.Auth.Logout)
	auth.GET("/me", s.svcs.Auth.Me, authmw.RequireAuth())

	// Short URL routes
	apiV1.POST("/shorten", s.svcs.ShortURL.CreateShortURL)
	apiV1.GET("/my-urls", s.svcs.ShortURL.GetMyURLs)
	apiV1.PATCH("/urls/:code/toggle", s.svcs.ShortURL.ToggleURLActive)
	apiV1.DELETE("/urls/:code", s.svcs.ShortURL.DeleteURL)

	// Redirect route at root level: example.com/:code (no rate limit for fast redirects)
	s.e.GET("/:code", s.svcs.ShortURL.GetOriginalURL)
}
