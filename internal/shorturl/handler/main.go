package handler

import (
	analyticsservice "url_shortner_backend/internal/analytics/service"
	"url_shortner_backend/internal/shorturl/service"
)

type ShortURLHandler struct {
	Svc          service.ShortURLSvc
	FrontendURL  string
	AnalyticsSvc analyticsservice.AnalyticsSvc
}

func NewShortURLHandler(svc service.ShortURLSvc, frontendURL string) *ShortURLHandler {
	return &ShortURLHandler{
		Svc:         svc,
		FrontendURL: frontendURL,
	}
}
