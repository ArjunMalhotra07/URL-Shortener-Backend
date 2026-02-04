package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

type GetSummaryRes struct {
	TotalClicks  int64                   `json:"total_clicks"`
	UniqueClicks int64                   `json:"unique_clicks"`
	BotClicks    int64                   `json:"bot_clicks"`
	TopReferrers []service.ReferrerStats `json:"top_referrers"`
}

func (h *AnalyticsHandler) GetSummary(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	timeRange := c.QueryParam("range")
	if timeRange == "" {
		timeRange = "7d" // default to 7 days
	}

	ownerType, ownerID, authExpired := getOwnerInfo(c)
	if authExpired {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	}

	output, err := h.Svc.GetSummary(c.Request().Context(), service.GetSummaryInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		TimeRange: timeRange,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetSummaryRes{
		TotalClicks:  output.TotalClicks,
		UniqueClicks: output.UniqueClicks,
		BotClicks:    output.BotClicks,
		TopReferrers: output.TopReferrers,
	})
}

// getOwnerInfo returns ownerType, ownerID, and whether auth has expired
func getOwnerInfo(c echo.Context) (string, string, bool) {
	// Check if auth token expired
	if c.Get("auth_expired") == true {
		return "", "", true
	}

	// Check if user is authenticated
	if userID := c.Get("user_id"); userID != nil && userID != "" {
		return "user", userID.(string), false
	}

	// Fall back to anonymous from cookie
	if cookie, err := c.Cookie("anon_id"); err == nil && cookie.Value != "" {
		// Parse the cookie value (format: "uuid_timestamp")
		anonID := cookie.Value
		if idx := len(anonID) - 11; idx > 0 && anonID[idx] == '_' {
			anonID = anonID[:idx]
		}
		return "anonymous", anonID, false
	}

	return "", "", false
}

func handleAnalyticsError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, service.ErrURLNotFound):
		return c.JSON(http.StatusNotFound, ErrorRes{Error: "URL not found"})
	case errors.Is(err, service.ErrURLNotOwned):
		return c.JSON(http.StatusForbidden, ErrorRes{Error: "you do not own this URL"})
	case errors.Is(err, service.ErrAnalyticsFetch):
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to fetch analytics"})
	default:
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
	}
}
