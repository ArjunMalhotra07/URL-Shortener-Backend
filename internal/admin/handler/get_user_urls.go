package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/admin/service"
)

func (h *AdminHandler) GetUserURLs(c echo.Context) error {
	userID := c.Param("id")
	if userID == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "user id is required"})
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	output, err := h.Svc.GetUserURLs(c.Request().Context(), service.GetUserURLsInput{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to get user urls"})
	}

	return c.JSON(http.StatusOK, output)
}
