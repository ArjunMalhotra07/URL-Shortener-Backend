package internal

import (
	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/shorturl/handler"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/internal/shorturl/service"
	"url_shortner_backend/pkg/logger"
)

type AppServices struct {
	ShortURL handler.ShortURLHandler
}
type AppServicesParams struct {
	Queries *db.Queries
	Logger  logger.Logger
}

func GetAppServices(p AppServicesParams) *AppServices {
	//Short Url
	shortURLRepo := repo.NewShortURLRepoImp(repo.ShortURLRepoParams{Queries: p.Queries})
	shortURLSvc := service.NewShortURLSvcImp(shortURLRepo, p.Logger)
	shortURLHandler := handler.NewShortURLHandler(shortURLSvc)

	return &AppServices{
		ShortURL: *shortURLHandler,
	}
}
