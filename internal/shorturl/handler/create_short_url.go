package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	AnonIDCookieName = "anon_id"
	CookieMaxAge     = 365 * 24 * 60 * 60 // 1 year in seconds
)

type CreateShortURLReq struct {
	LongURL string `json:"long_url"`
}

type CreateShortURLRes struct {
	Code      string `json:"code"`
	LongURL   string `json:"long_url"`
	OwnerType string `json:"owner_type"`
}

type ErrorRes struct {
	Error string `json:"error"`
}

func (h *ShortURLHandler) CreateShortURL(c echo.Context) error {
	var req CreateShortURLReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.LongURL == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "long_url is required"})
	}

	// Get or create anon_id from cookie
	anonID := getOrCreateAnonID(c)

	// TODO: Check if user is authenticated (JWT), use user_id instead
	// For now, always use anonymous
	ownerType := "anonymous"
	ownerID := anonID

	output, err := h.Svc.CreateShortURL(c.Request().Context(), service.CreateShortURLInput{
		LongURL:   req.LongURL,
		OwnerType: ownerType,
		OwnerID:   ownerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLCreation), errors.Is(err, service.ErrURLCodeUpdate):
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to create short url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	return c.JSON(http.StatusCreated, CreateShortURLRes{
		Code:      output.Code,
		LongURL:   output.LongURL,
		OwnerType: output.OwnerType,
	})
}

// getOrCreateAnonID gets existing anon_id from cookie or creates a new one
func getOrCreateAnonID(c echo.Context) string {
	cookie, err := c.Cookie(AnonIDCookieName)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Generate new UUID
	anonID := uuid.New().String()

	// Set cookie
	c.SetCookie(&http.Cookie{
		Name:     AnonIDCookieName,
		Value:    anonID,
		Path:     "/",
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Secure:   true, // Set to false for local dev if not using HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return anonID
}
