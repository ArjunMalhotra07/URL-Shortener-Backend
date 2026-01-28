package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

func (h *AuthHandler) Refresh(c echo.Context) error {
	cookie, err := c.Cookie(RefreshTokenCookie)
	if err != nil || cookie.Value == "" {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "refresh token required"})
	}

	output, err := h.Svc.Refresh(c.Request().Context(), cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidToken), errors.Is(err, service.ErrTokenExpired):
			clearAuthCookies(c)
			return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "invalid or expired refresh token"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "refresh failed"})
		}
	}

	setAuthCookies(c, output)

	return c.JSON(http.StatusOK, AuthRes{
		UserID: output.UserID,
		Email:  output.Email,
	})
}
