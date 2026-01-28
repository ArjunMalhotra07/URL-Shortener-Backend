package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *AuthHandler) Logout(c echo.Context) error {
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err == nil && cookie.Value != "" {
		_ = h.Svc.Logout(c.Request().Context(), cookie.Value)
	}

	clearAuthCookies(c)

	return c.JSON(http.StatusOK, map[string]string{"message": "logged out"})
}
