package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal"
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/httpserver"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
	"url_shortner_backend/pkg/migrate"
	"url_shortner_backend/pkg/postgres"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("couldn't start server no config found")
	}
	// Load logger
	logr := logger.NewZeroLogger()

	// Run migrations
	migrate.RunMigrations(cfg.DBDSN)

	// Connect to DB
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, postgres.Params{
		DSN:             cfg.DBDSN,
		MaxConns:        cfg.DBMaxConns,
		MaxConnLifetime: time.Duration(cfg.DBMaxConnLifetimeSeconds) * time.Second,
		MaxConnIdleTime: time.Duration(cfg.DBMaxConnIdleTimeSeconds) * time.Second,
	})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()
	queries := db.New(pool)

	// JWT Manager
	jwtMgr := jwt.NewJWTManager(jwt.Config{
		Secret:             cfg.JWTSecret,
		AccessTokenExpiry:  time.Duration(cfg.JWTAccessTokenExpiryMins) * time.Minute,
		RefreshTokenExpiry: time.Duration(cfg.JWTRefreshTokenExpiryDays) * 24 * time.Hour,
	})

	// Get services
	svcs := internal.GetAppServices(internal.AppServicesParams{
		Queries:     queries,
		Logger:      logr,
		JWT:         jwtMgr,
		FrontendURL: cfg.FrontendURL,
	})

	// Server
	svr := httpserver.NewEchoServer(svcs, jwtMgr)

	go func() {
		logr.Info("Starting server", "addr", cfg.ServerPort)
		if err := svr.Start(cfg.ServerPort); err != nil {
			log.Fatal("server stopped: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := svr.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}
