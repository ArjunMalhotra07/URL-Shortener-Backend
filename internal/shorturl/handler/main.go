package handler

import "url_shortner_backend/internal/shorturl/service"

type ShortURLHandler struct {
	Svc service.ShortURLSvc
}

func GetShortURLHandler(svc service.ShortURLSvc) *ShortURLHandler {
	return &ShortURLHandler{Svc: svc}
}
