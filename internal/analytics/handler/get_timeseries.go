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

	timeRange, errMsg := parseTimeRangeParams(c)
	if errMsg != "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: errMsg})
	}

	interval := c.QueryParam("interval")
	if interval == "" {
		// Auto-select based on time range
		rangeParam := c.QueryParam("range")
		if rangeParam == "24h" {
			interval = "hour"
		} else {
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
		Start:     timeRange.Start,
		End:       timeRange.End,
		Interval:  interval,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetTimeseriesRes{
		Data: output.Data,
	})
}
