package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/httpserver"
	"url_shortner_backend/pkg/logger"
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

	// Connect to DB
	ctx := context.Background()
	pool, err := postgres.NewPool(ctx, postgres.Params{
		DSN:             cfg.DBDSN,
		MinConns:        cfg.DBMinConns,
		MaxConns:        cfg.DBMaxConns,
		MaxConnLifetime: time.Duration(cfg.DBMaxConnLifetimeSeconds) * time.Second,
	})
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}
	defer pool.Close()
	// Get services

	// Server
	svr := httpserver.EchoServer()

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
