package handler

import (
	"net/http"
	"strconv"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

type GetClicksRes struct {
	Clicks     []service.ClickRecord `json:"clicks"`
	TotalCount int64                 `json:"total_count"`
}

func (h *AnalyticsHandler) GetClicks(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	limit := int32(50)
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = int32(l)
		}
	}

	offset := int32(0)
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = int32(o)
		}
	}

	ownerType, ownerID, authExpired := getOwnerInfo(c)
	if authExpired {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	}

	output, err := h.Svc.GetClicks(c.Request().Context(), service.GetClicksInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetClicksRes{
		Clicks:     output.Clicks,
		TotalCount: output.TotalCount,
	})
}
