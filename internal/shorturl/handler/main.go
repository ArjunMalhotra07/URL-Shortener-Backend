package handler

import "url_shortner_backend/internal/shorturl/service"

type ShortURLHandler struct {
	Svc         service.ShortURLSvc
	FrontendURL string
}

func NewShortURLHandler(svc service.ShortURLSvc, frontendURL string) *ShortURLHandler {
	return &ShortURLHandler{
		Svc:         svc,
		FrontendURL: frontendURL,
	}
}
