package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AdminHandler) GetStats(c echo.Context) error {
	output, err := h.Svc.GetPlatformStats(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to get platform stats"})
	}

	return c.JSON(http.StatusOK, output)
}
