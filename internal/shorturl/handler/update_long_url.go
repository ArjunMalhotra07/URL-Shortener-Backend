package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

type UpdateLongURLReq struct {
	LongURL string `json:"long_url"`
}

func (h *ShortURLHandler) UpdateLongURL(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	var req UpdateLongURLReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.LongURL == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "long_url is required"})
	}

	// Get owner info from context
	var ownerType, ownerID string
	if userID := c.Get("user_id"); userID != nil && userID != "" {
		ownerType = "user"
		ownerID = userID.(string)
	} else if c.Get("auth_expired") == true {
		// User had a token but it expired - return 401 so frontend can refresh
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "token expired"})
	} else {
		// Check for anon_id cookie
		cookie, err := c.Cookie(AnonIDCookieName)
		if err != nil || cookie.Value == "" {
			return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "authentication required"})
		}
		ownerType = "anonymous"
		anonID, _, valid := parseAnonCookie(cookie.Value)
		if !valid {
			return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "invalid session"})
		}
		ownerID = anonID
	}

	err := h.Svc.UpdateLongURL(c.Request().Context(), service.UpdateLongURLInput{
		Code:      code,
		LongURL:   req.LongURL,
		OwnerType: ownerType,
		OwnerID:   ownerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid url"})
		case errors.Is(err, service.ErrURLNotOwned):
			return c.JSON(http.StatusForbidden, ErrorRes{Error: "you don't own this url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to update url"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "url updated successfully"})
}
