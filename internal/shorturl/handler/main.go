package handler

import (
	"context"

	analyticsservice "url_shortner_backend/internal/analytics/service"
	"url_shortner_backend/internal/shorturl/service"
)

type BlockChecker interface {
	IsOwnerBlocked(ctx context.Context, code string) (bool, error)
}

type ShortURLHandler struct {
	Svc          service.ShortURLSvc
	FrontendURL  string
	AnalyticsSvc analyticsservice.AnalyticsSvc
	BlockChecker BlockChecker
}

func NewShortURLHandler(svc service.ShortURLSvc, frontendURL string) *ShortURLHandler {
	return &ShortURLHandler{
		Svc:         svc,
		FrontendURL: frontendURL,
	}
}
