package handler

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func isProd() bool {
	return os.Getenv("APP_ENV") != "dev"
}

const (
	AnonIDCookieName   = "anon_id"
	CookieMaxAge       = 365 * 24 * 60 * 60 // 1 year in seconds
	CookieRefreshAfter = 335 * 24 * 60 * 60 // Refresh if older than 11 months (335 days)
)

type CreateShortURLReq struct {
	LongURL   string  `json:"long_url"`
	ExpiresAt *string `json:"expires_at,omitempty"`
}

type CreateShortURLRes struct {
	Code                 string  `json:"code"`
	LongURL              string  `json:"long_url"`
	OwnerType            string  `json:"owner_type"`
	IsActive             bool    `json:"is_active"`
	ExpiresAt            *string `json:"expires_at,omitempty"`
	CreatedAt            string  `json:"created_at"`
	URLsCreatedThisMonth int64   `json:"urls_created_this_month"`
	URLsLimit            int     `json:"urls_limit"`
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

	// Check if user is authenticated
	var ownerType, ownerID string
	var expiresAt *time.Time
	if userID := c.Get("user_id"); userID != nil && userID != "" {
		ownerType = "user"
		ownerID = userID.(string)
		if req.ExpiresAt != nil {
			parsed, err := time.Parse(time.RFC3339, *req.ExpiresAt)
			if err != nil {
				return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid expires_at format, use RFC3339"})
			}
			if parsed.Before(time.Now()) {
				return c.JSON(http.StatusBadRequest, ErrorRes{Error: "expires_at must be in the future"})
			}
			expiresAt = &parsed
		}
	} else {
		// Fall back to anonymous
		ownerType = "anonymous"
		ownerID = getOrCreateAnonID(c)
	}

	output, err := h.Svc.CreateShortURL(c.Request().Context(), service.CreateShortURLInput{
		LongURL:   req.LongURL,
		OwnerType: ownerType,
		OwnerID:   ownerID,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrMonthlyQuotaExceeded):
			return c.JSON(http.StatusTooManyRequests, ErrorRes{Error: "Monthly quota exceeded"})
		case errors.Is(err, service.ErrURLCreation), errors.Is(err, service.ErrURLCodeUpdate):
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to create short url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	var expiresAtStr *string
	if output.ExpiresAt != nil {
		formatted := output.ExpiresAt.UTC().Format(time.RFC3339)
		expiresAtStr = &formatted
	}

	return c.JSON(http.StatusCreated, CreateShortURLRes{
		Code:                 output.Code,
		LongURL:              output.LongURL,
		OwnerType:            output.OwnerType,
		IsActive:             output.IsActive,
		ExpiresAt:            expiresAtStr,
		CreatedAt:            output.CreatedAt.UTC().Format(time.RFC3339),
		URLsCreatedThisMonth: output.URLsCreatedThisMonth,
		URLsLimit:            output.URLsLimit,
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

	sameSite := http.SameSiteLaxMode
	secure := false
	if isProd() {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	c.SetCookie(&http.Cookie{
		Name:     AnonIDCookieName,
		Value:    cookieValue,
		Path:     "/",
		MaxAge:   CookieMaxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}
