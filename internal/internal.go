package internal

import (
	db "url_shortner_backend/db/output"
	"url_shortner_backend/internal/shorturl/handler"
	"url_shortner_backend/internal/shorturl/repo"
	"url_shortner_backend/internal/shorturl/service"
)

type AppServices struct {
	ShortURL handler.ShortURLHandler
}
type AppServicesParams struct {
	Queries *db.Queries
}

func GetAppServices(p AppServicesParams) *AppServices {
	//Short Url
	shortURLRepo := repo.NewShortURLRepoImp(repo.ShortURLRepoParams{Queries: p.Queries})
	shortURLSvc := service.NewShortURLSvcImp(shortURLRepo)
	shortURLHandler := handler.NewShortURLHandler(shortURLSvc)

	return &AppServices{
		ShortURL: *shortURLHandler,
	}
}
