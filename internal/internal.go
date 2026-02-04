package internal

import (
	db "url_shortner_backend/db/output"
	analyticshandler "url_shortner_backend/internal/analytics/handler"
	analyticsrepo "url_shortner_backend/internal/analytics/repo"
	analyticsservice "url_shortner_backend/internal/analytics/service"
	authhandler "url_shortner_backend/internal/auth/handler"
	authrepo "url_shortner_backend/internal/auth/repo"
	authservice "url_shortner_backend/internal/auth/service"
	"url_shortner_backend/internal/shorturl/handler"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/internal/shorturl/service"
	"url_shortner_backend/pkg/config"
	"url_shortner_backend/pkg/geoip"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
	"url_shortner_backend/pkg/redis"
)

type AppServices struct {
	ShortURL  handler.ShortURLHandler
	Auth      authhandler.AuthHandler
	Analytics analyticshandler.AnalyticsHandler
}

// AnalyticsSvc exposes the analytics service for use in redirect recording
type AnalyticsSvc = analyticsservice.AnalyticsSvc

type AppServicesParams struct {
	Queries *db.Queries
	Logger  logger.Logger
	JWT     *jwt.JWTManager
	Cfg     *config.Config
	Redis   *redis.RedisClient
	GeoIP   geoip.GeoIPLookup
}

func GetAppServices(p AppServicesParams) *AppServices {
	// Short URL
	shortURLRepo := repo.NewShortURLRepoImp(repo.ShortURLRepoParams{Queries: p.Queries})
	shortURLSvc := service.NewShortURLSvcImp(shortURLRepo, p.Logger, p.Cfg, p.Redis)
	shortURLHandler := handler.NewShortURLHandler(shortURLSvc, p.Cfg.FrontendURL)

	// Auth
	authRepo := authrepo.NewAuthRepoImp(authrepo.AuthRepoParams{Queries: p.Queries})
	authSvc := authservice.NewAuthSvcImp(authRepo, shortURLSvc, p.JWT, p.Logger)
	authHandler := authhandler.NewAuthHandler(authSvc)

	// Analytics
	// Toggle: Use dummy repo for frontend testing (set USE_DUMMY_ANALYTICS=true)
	var analyticsRepo analyticsrepo.AnalyticsRepository
	if p.Cfg.UseDummyAnalytics {
		analyticsRepo = analyticsrepo.NewAnalyticsRepoDummy()
	} else {
		analyticsRepo = analyticsrepo.NewAnalyticsRepoImp(analyticsrepo.AnalyticsRepoParams{Queries: p.Queries})
	}
	analyticsSvc := analyticsservice.NewAnalyticsSvcImp(analyticsservice.AnalyticsSvcParams{
		Logger:       p.Logger,
		Repo:         analyticsRepo,
		ShortURLRepo: shortURLRepo,
		GeoIP:        p.GeoIP,
		Redis:        p.Redis,
	})
	analyticsHandler := analyticshandler.NewAnalyticsHandler(analyticsSvc)

	// Set analytics service on shortURL handler for click recording
	shortURLHandler.AnalyticsSvc = analyticsSvc

	return &AppServices{
		ShortURL:  *shortURLHandler,
		Auth:      *authHandler,
		Analytics: *analyticsHandler,
	}
}
