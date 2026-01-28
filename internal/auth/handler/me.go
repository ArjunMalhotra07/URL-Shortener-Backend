package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id")
	email := c.Get("email")

	if userID == nil || userID == "" {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "not authenticated"})
	}

	return c.JSON(http.StatusOK, AuthRes{
		UserID: userID.(string),
		Email:  email.(string),
	})
}
