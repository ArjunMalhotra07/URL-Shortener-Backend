package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/admin/service"
)

func (h *AdminHandler) ListUsers(c echo.Context) error {
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	output, err := h.Svc.ListUsers(c.Request().Context(), service.ListUsersInput{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to list users"})
	}

	return c.JSON(http.StatusOK, output)
}
