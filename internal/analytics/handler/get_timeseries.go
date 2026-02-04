package handler

import (
	"net/http"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

type GetTimeseriesRes struct {
	Data []service.TimeseriesPoint `json:"data"`
}

func (h *AnalyticsHandler) GetTimeseries(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	timeRange := c.QueryParam("range")
	if timeRange == "" {
		timeRange = "7d"
	}

	interval := c.QueryParam("interval")
	if interval == "" {
		// Auto-select based on time range
		switch timeRange {
		case "24h":
			interval = "hour"
		default:
			interval = "day"
		}
	}

	ownerType, ownerID, authExpired := getOwnerInfo(c)
	if authExpired {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	}

	output, err := h.Svc.GetTimeseries(c.Request().Context(), service.GetTimeseriesInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		TimeRange: timeRange,
		Interval:  interval,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetTimeseriesRes{
		Data: output.Data,
	})
}
