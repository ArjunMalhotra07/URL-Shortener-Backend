package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

func isProd() bool {
	return os.Getenv("APP_ENV") != "dev"
}

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
	UserID  string `json:"user_id"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

type ErrorRes struct {
	Error string `json:"error"`
}

func setAuthCookies(c echo.Context, output service.AuthOutput) {
	sameSite := http.SameSiteLaxMode
	secure := false
	if isProd() {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	c.SetCookie(&http.Cookie{
		Name:     AccessTokenCookie,
		Value:    output.AccessToken,
		Path:     "/",
		Expires:  output.AccessExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    output.RefreshToken,
		Path:     "/api/v1/auth",
		Expires:  output.RefreshExpiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func clearAuthCookies(c echo.Context) {
	sameSite := http.SameSiteLaxMode
	secure := false
	if isProd() {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	c.SetCookie(&http.Cookie{
		Name:     AccessTokenCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})

	c.SetCookie(&http.Cookie{
		Name:     RefreshTokenCookie,
		Value:    "",
		Path:     "/api/v1/auth",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

func clearAnonCookie(c echo.Context) {
	sameSite := http.SameSiteLaxMode
	secure := false
	if isProd() {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	c.SetCookie(&http.Cookie{
		Name:     AnonIDCookie,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
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
