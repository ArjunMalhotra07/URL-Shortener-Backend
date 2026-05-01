package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/pkg/jwt"
)

const AdminTokenCookie = "admin_token"

// AdminAuthMiddleware validates the admin JWT from the admin_token cookie.
// Sets "is_admin" in context if valid.
func AdminAuthMiddleware(jwtMgr *jwt.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(AdminTokenCookie)
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "admin authentication required",
				})
			}

			claims, err := jwtMgr.ValidateAccessToken(cookie.Value)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired admin token",
				})
			}

			if claims.UserID != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "not an admin token",
				})
			}

			c.Set("is_admin", true)
			return next(c)
		}
	}
}
