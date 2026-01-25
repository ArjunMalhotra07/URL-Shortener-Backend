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
)

func main() {
	//1. get config
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("couldn't start server no config found")
	}
	//2. get logger
	logger := logger.NewZeroLogger()
	//3. get services
	// svcs := services.NewServices(models.Params{
	// 	Config: cfg,
	// 	Logger: logger,
	// })
	//4. get server
	httpServer := httpserver.NewEchoServer()
	// 5. Init app
	srv := NewServer(ServerParams{
		Config: cfg,
		Http:   httpServer,
		Logger: logger,
		// Services: svcs,
	})
	//! Start server
	go func() {
		logger.Info("Starting server on port " + cfg.ServerPort)
		if err := srv.Http.Start(cfg.ServerPort); err != nil {
			log.Fatal("shutting down the server")
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Http.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
