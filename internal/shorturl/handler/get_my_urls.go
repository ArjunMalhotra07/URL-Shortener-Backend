package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

type URLItem struct {
	Code      string  `json:"code"`
	LongURL   string  `json:"long_url"`
	IsActive  bool    `json:"is_active"`
	ExpiresAt *string `json:"expires_at,omitempty"`
	IsExpired bool    `json:"is_expired"`
	CreatedAt string  `json:"created_at"`
	Name      *string `json:"name,omitempty"`
}

type GetMyURLsRes struct {
	URLs   []URLItem `json:"urls"`
	Total  int64     `json:"total"`
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
}

func (h *ShortURLHandler) GetMyURLs(c echo.Context) error {
	// Parse pagination params
	limit := int32(10)
	offset := int32(0)
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = int32(parsed)
		}
	}
	if o := c.QueryParam("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = int32(parsed)
		}
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
		// Fall back to anonymous
		ownerType = "anonymous"
		ownerID = getOrCreateAnonID(c)
	}

	output, err := h.Svc.GetMyURLs(c.Request().Context(), service.GetMyURLsInput{
		OwnerType: ownerType,
		OwnerID:   ownerID,
		Limit:     limit,
		Offset:    offset,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOwner):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLFetch):
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to fetch urls"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	urls := make([]URLItem, len(output.URLs))
	for i, u := range output.URLs {
		var expiresAt *string
		if u.ExpiresAt != nil {
			formatted := u.ExpiresAt.UTC().Format(time.RFC3339)
			expiresAt = &formatted
		}
		urls[i] = URLItem{
			Code:      u.Code,
			LongURL:   u.LongURL,
			IsActive:  u.IsActive,
			ExpiresAt: expiresAt,
			IsExpired: u.IsExpired,
			CreatedAt: u.CreatedAt.UTC().Format(time.RFC3339),
			Name:      u.Name,
		}
	}

	return c.JSON(http.StatusOK, GetMyURLsRes{
		URLs:   urls,
		Total:  output.Total,
		Limit:  limit,
		Offset: offset,
	})
}
