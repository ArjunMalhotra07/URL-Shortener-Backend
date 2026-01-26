package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	AnonIDCookieName   = "anon_id"
	CookieMaxAge       = 365 * 24 * 60 * 60 // 1 year in seconds
	CookieRefreshAfter = 335 * 24 * 60 * 60 // Refresh if older than 11 months (335 days)
)

type CreateShortURLReq struct {
	LongURL string `json:"long_url"`
}

type CreateShortURLRes struct {
	Code      string `json:"code"`
	LongURL   string `json:"long_url"`
	OwnerType string `json:"owner_type"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
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
		IsActive:  output.IsActive,
		CreatedAt: output.CreatedAt.Format("2006-01-02T15:04:05Z"),
	})
}

func getOrCreateAnonID(c echo.Context) string {
	cookie, err := c.Cookie(AnonIDCookieName)
	if err == nil && cookie.Value != "" {
		anonID, createdAt, valid := parseAnonCookie(cookie.Value)
		if valid {
			if time.Since(createdAt).Seconds() > CookieRefreshAfter {
				setAnonCookie(c, anonID)
			}
			return anonID
		}
	}
	anonID := uuid.New().String()
	setAnonCookie(c, anonID)

	return anonID
}

func parseAnonCookie(value string) (string, time.Time, bool) {
	parts := strings.Split(value, "_")
	if len(parts) != 2 {
		return "", time.Time{}, false
	}

	anonID := parts[0]
	timestamp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", time.Time{}, false
	}

	return anonID, time.Unix(timestamp, 0), true
}

func setAnonCookie(c echo.Context, anonID string) {
	cookieValue := fmt.Sprintf("%s_%d", anonID, time.Now().Unix())

	c.SetCookie(&http.Cookie{
		Name:     AnonIDCookieName,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}
