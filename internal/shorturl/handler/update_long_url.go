package handler

import (
	"errors"
	"net/http"
	"time"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

type UpdateLongURLReq struct {
	LongURL   string  `json:"long_url"`
	ExpiresAt *string `json:"expires_at"` // ISO 8601 timestamp, null to remove expiry
	Name      *string `json:"name,omitempty"`
}

type UpdateLongURLRes struct {
	Code      string  `json:"code"`
	LongURL   string  `json:"long_url"`
	IsActive  bool    `json:"is_active"`
	ExpiresAt *string `json:"expires_at,omitempty"`
	IsExpired bool    `json:"is_expired"`
	CreatedAt string  `json:"created_at"`
	Name      *string `json:"name,omitempty"`
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

	// Parse expires_at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			// Empty string means remove expiry
			zeroTime := time.Time{}
			expiresAt = &zeroTime
		} else {
			parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid expires_at format, use ISO 8601 (RFC3339)"})
			}
			expiresAt = &parsed
		}
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

	output, err := h.Svc.UpdateLongURL(c.Request().Context(), service.UpdateLongURLInput{
		Code:      code,
		LongURL:   req.LongURL,
		ExpiresAt: expiresAt,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Name:      req.Name,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid url"})
		case errors.Is(err, service.ErrInvalidExpiresAt):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: "expires_at must be in the future"})
		case errors.Is(err, service.ErrURLNotOwned):
			return c.JSON(http.StatusForbidden, ErrorRes{Error: "you don't own this url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to update url"})
		}
	}

	// Build response
	res := UpdateLongURLRes{
		Code:      output.Code,
		LongURL:   output.LongURL,
		IsActive:  output.IsActive,
		IsExpired: output.IsExpired,
		CreatedAt: output.CreatedAt.UTC().Format(time.RFC3339),
		Name:      output.Name,
	}

	if output.ExpiresAt != nil {
		formatted := output.ExpiresAt.UTC().Format(time.RFC3339)
		res.ExpiresAt = &formatted
	}

	return c.JSON(http.StatusOK, res)
}
