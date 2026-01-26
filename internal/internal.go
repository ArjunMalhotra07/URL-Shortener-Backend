package internal

import "url_shortner_backend/internal/shorturl/handler"

type AppServices struct {
	ShortURL handler.ShortURLHandler
}
type AppServicesParams struct {
}

func GetAppServices() *AppServices {

	return &AppServices{}
}
