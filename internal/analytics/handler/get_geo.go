package handler

import (
	"net/http"
	"strconv"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

type GetGeoRes struct {
	Countries []service.CountryStats `json:"countries"`
	Cities    []service.CityStats    `json:"cities"`
}

func (h *AnalyticsHandler) GetGeo(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	timeRange, errMsg := parseTimeRangeParams(c)
	if errMsg != "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: errMsg})
	}

	limit := int32(10)
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 50 {
			limit = int32(l)
		}
	}

	ownerType, ownerID, authExpired := getOwnerInfo(c)
	if authExpired {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	}

	output, err := h.Svc.GetGeoBreakdown(c.Request().Context(), service.GetGeoInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Start:     timeRange.Start,
		End:       timeRange.End,
		Limit:     limit,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetGeoRes{
		Countries: output.Countries,
		Cities:    output.Cities,
	})
}
