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
	"url_shortner_backend/pkg/geoip"
	"url_shortner_backend/pkg/httpserver"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
	"url_shortner_backend/pkg/migrate"
	"url_shortner_backend/pkg/postgres"
	"url_shortner_backend/pkg/redis"
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

	// Redis client (optional - analytics still works without it)
	var redisClient *redis.RedisClient
	if cfg.RedisAddr != "" {
		redisClient, err = redis.NewRedisClient(redis.RedisConfig{
			Addr:     cfg.RedisAddr,
			Password: cfg.RedisPassword,
			DB:       cfg.RedisDB,
		})
		if err != nil {
			logr.Info("redis connection failed, analytics caching disabled", "error", err)
		} else {
			defer redisClient.Close()
			logr.Info("connected to redis", "addr", cfg.RedisAddr)
		}
	}

	// GeoIP service (optional - analytics still works without it)
	var geoIPSvc geoip.GeoIPLookup
	if cfg.GeoIPDBPath != "" {
		realGeoIP, err := geoip.NewGeoIPService(cfg.GeoIPDBPath)
		if err != nil {
			logr.Info("GeoIP database not loaded, geo lookups disabled", "error", err, "path", cfg.GeoIPDBPath)
			geoIPSvc = geoip.NewNullGeoIPService()
		} else {
			defer realGeoIP.Close()
			geoIPSvc = realGeoIP
			logr.Info("loaded GeoIP database", "path", cfg.GeoIPDBPath)
		}
	} else {
		geoIPSvc = geoip.NewNullGeoIPService()
	}

	// Get services
	svcs := internal.GetAppServices(internal.AppServicesParams{
		Queries: queries,
		Logger:  logr,
		JWT:     jwtMgr,
		Cfg:     &cfg,
		Redis:   redisClient,
		GeoIP:   geoIPSvc,
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
