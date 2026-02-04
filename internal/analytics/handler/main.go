package handler

import "url_shortner_backend/internal/analytics/service"

type AnalyticsHandler struct {
	Svc service.AnalyticsSvc
}

func NewAnalyticsHandler(svc service.AnalyticsSvc) *AnalyticsHandler {
	return &AnalyticsHandler{
		Svc: svc,
	}
}

type ErrorRes struct {
	Error string `json:"error"`
}
