package internal

import (
	db "url_shortner_backend/db/output"
	authhandler "url_shortner_backend/internal/auth/handler"
	authrepo "url_shortner_backend/internal/auth/repo"
	authservice "url_shortner_backend/internal/auth/service"
	"url_shortner_backend/internal/shorturl/handler"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/internal/shorturl/service"
	"url_shortner_backend/pkg/jwt"
	"url_shortner_backend/pkg/logger"
)

type AppServices struct {
	ShortURL handler.ShortURLHandler
	Auth     authhandler.AuthHandler
}

type AppServicesParams struct {
	Queries       *db.Queries
	Logger        logger.Logger
	JWT           *jwt.JWTManager
	FrontendURL   string
	DailyURLQuota int
}

func GetAppServices(p AppServicesParams) *AppServices {
	// Short URL
	shortURLRepo := repo.NewShortURLRepoImp(repo.ShortURLRepoParams{Queries: p.Queries})
	shortURLSvc := service.NewShortURLSvcImp(shortURLRepo, p.Logger, p.DailyURLQuota)
	shortURLHandler := handler.NewShortURLHandler(shortURLSvc, p.FrontendURL)

	// Auth
	authRepo := authrepo.NewAuthRepoImp(authrepo.AuthRepoParams{Queries: p.Queries})
	authSvc := authservice.NewAuthSvcImp(authRepo, shortURLRepo, p.JWT, p.Logger, p.DailyURLQuota)
	authHandler := authhandler.NewAuthHandler(authSvc)

	return &AppServices{
		ShortURL: *shortURLHandler,
		Auth:     *authHandler,
	}
}
