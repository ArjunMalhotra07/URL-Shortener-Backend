package handler

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/admin/service"
)

const AdminTokenCookie = "admin_token"

type AdminHandler struct {
	Svc service.AdminService
}

func NewAdminHandler(svc service.AdminService) *AdminHandler {
	return &AdminHandler{Svc: svc}
}

type ErrorRes struct {
	Error string `json:"error"`
}

func isProd() bool {
	return os.Getenv("APP_ENV") != "dev"
}

func setAdminCookie(c echo.Context, token string, expiresAt time.Time) {
	sameSite := http.SameSiteLaxMode
	secure := false
	if isProd() {
		sameSite = http.SameSiteNoneMode
		secure = true
	}

	c.SetCookie(&http.Cookie{
		Name:     AdminTokenCookie,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}
