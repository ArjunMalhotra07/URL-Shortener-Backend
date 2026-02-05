package handler

import (
	"time"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

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

// TimeRange represents parsed time range for analytics queries
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// parseTimeRangeParams parses time range from query parameters
// Priority: custom range (start/end) > preset range > default (7d)
// Returns TimeRange and error message if invalid
func parseTimeRangeParams(c echo.Context) (TimeRange, string) {
	now := time.Now()

	// Check for custom range first
	startStr := c.QueryParam("start")
	endStr := c.QueryParam("end")

	if startStr != "" && endStr != "" {
		start, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			return TimeRange{}, "invalid start date format, use ISO 8601 (RFC3339)"
		}
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			return TimeRange{}, "invalid end date format, use ISO 8601 (RFC3339)"
		}
		if end.Before(start) {
			return TimeRange{}, "end date must be after start date"
		}
		if end.After(now) {
			end = now
		}
		return TimeRange{Start: start, End: end}, ""
	}

	// Fall back to preset range
	rangeParam := c.QueryParam("range")
	if rangeParam == "" {
		rangeParam = "7d" // default
	}

	var start time.Time
	switch rangeParam {
	case "24h":
		start = now.Add(-24 * time.Hour)
	case "7d":
		start = now.Add(-7 * 24 * time.Hour)
	case "30d":
		start = now.Add(-30 * 24 * time.Hour)
	case "90d":
		start = now.Add(-90 * 24 * time.Hour)
	case "all":
		start = time.Time{} // zero time means all
	default:
		start = now.Add(-7 * 24 * time.Hour) // default to 7d for unknown values
	}

	return TimeRange{Start: start, End: now}, ""
}
