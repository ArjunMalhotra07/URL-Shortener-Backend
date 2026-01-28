package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

const (
	AccessTokenCookie  = "access_token"
	RefreshTokenCookie = "refresh_token"
	AnonIDCookie       = "anon_id"
)

type AuthHandler struct {
	Svc service.AuthService
}

func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{Svc: svc}
}

type AuthRes struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type ErrorRes struct {
	Error string `json:"error"`
}

func setAuthCookies(c echo.Context, output service.AuthOutput) {
	c.SetCookie(&http.Cookie{
		Name:     AccessTokenCookie,
		Value:    output.AccessToken,
		Path:     "/",
		Expires:  output.AccessExpiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    output.RefreshToken,
		Path:     "/api/v1/auth", // Only sent to auth endpoints
		Expires:  output.RefreshExpiresAt,
		HttpOnly: true,
		Secure:   false, // Set to true in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAuthCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

func clearAnonCookie(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     AnonIDCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
}

// parseAnonID extracts the UUID from the anon_id cookie value (format: uuid_timestamp)
func parseAnonID(value string) string {
	for i := len(value) - 1; i >= 0; i-- {
		if value[i] == '_' {
			return value[:i]
		}
	}
	return value
}
