package handler

import (
	"context"
	"errors"
	"net/http"

	analyticsservice "url_shortner_backend/internal/analytics/service"
	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

func (h *ShortURLHandler) GetOriginalURL(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=invalid")
	}

	output, err := h.Svc.GetLongURL(c.Request().Context(), service.GetLongURLInput{
		Code: code,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCode):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=invalid")
		case errors.Is(err, service.ErrURLNotFound):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=not_found")
		case errors.Is(err, service.ErrURLExpired):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=expired")
		case errors.Is(err, service.ErrURLInactive):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=inactive")
		default:
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=default")
		}
	}

	// Record click asynchronously (don't block redirect)
	if h.AnalyticsSvc != nil {
		go h.recordClick(output.ID, c)
	}

	return c.Redirect(http.StatusFound, output.LongURL)
}

func (h *ShortURLHandler) recordClick(shortURLID int64, c echo.Context) {
	// Use background context since the request context will be cancelled after redirect
	ctx := context.Background()

	// Extract request metadata
	ipAddress := c.RealIP()
	userAgent := c.Request().Header.Get("User-Agent")
	referrer := c.Request().Header.Get("Referer")
	utmSource := c.QueryParam("utm_source")
	utmMedium := c.QueryParam("utm_medium")
	utmCampaign := c.QueryParam("utm_campaign")

	_ = h.AnalyticsSvc.RecordClick(ctx, analyticsservice.RecordClickInput{
		ShortURLID:  shortURLID,
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Referrer:    referrer,
		UTMSource:   utmSource,
		UTMMedium:   utmMedium,
		UTMCampaign: utmCampaign,
	})
}
