package handler

import "url_shortner_backend/internal/shorturl/service"

type ShortURLHandler struct {
	Svc service.ShortURLSvc
}

func NewShortURLHandler(svc service.ShortURLSvc) *ShortURLHandler {
	return &ShortURLHandler{Svc: svc}
}
