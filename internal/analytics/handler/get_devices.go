package handler

import (
	"net/http"

	"url_shortner_backend/internal/analytics/service"

	"github.com/labstack/echo/v4"
)

type GetDevicesRes struct {
	DeviceTypes []service.DeviceTypeStats `json:"device_types"`
	Browsers    []service.BrowserStats    `json:"browsers"`
	OS          []service.OSStats         `json:"os"`
}

func (h *AnalyticsHandler) GetDevices(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	timeRange, errMsg := parseTimeRangeParams(c)
	if errMsg != "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: errMsg})
	}

	ownerType, ownerID, authExpired := getOwnerInfo(c)
	if authExpired {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	}

	output, err := h.Svc.GetDeviceBreakdown(c.Request().Context(), service.GetDeviceInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Start:     timeRange.Start,
		End:       timeRange.End,
	})
	if err != nil {
		return handleAnalyticsError(c, err)
	}

	return c.JSON(http.StatusOK, GetDevicesRes{
		DeviceTypes: output.DeviceTypes,
		Browsers:    output.Browsers,
		OS:          output.OS,
	})
}
