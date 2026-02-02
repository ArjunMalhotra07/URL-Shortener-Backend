package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/shorturl/service"
)

func (h *ShortURLHandler) ToggleURLActive(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	// Check if user is authenticated
	var ownerType, ownerID string
	if userID := c.Get("user_id"); userID != nil && userID != "" {
		ownerType = "user"
		ownerID = userID.(string)
	} else if c.Get("auth_expired") == true {
		// User had a token but it expired - return 401 so frontend can refresh
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	} else {
		ownerType = "anonymous"
		ownerID = getOrCreateAnonID(c)
	}

	err := h.Svc.ToggleURLActive(c.Request().Context(), service.ToggleURLInput{
		Code:      code,
		OwnerType: ownerType,
		OwnerID:   ownerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCode):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLNotOwned):
			return c.JSON(http.StatusForbidden, ErrorRes{Error: "you don't own this url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to toggle url"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}
