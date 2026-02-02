package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/pkg/jwt"
)

const (
	AccessTokenCookie = "access_token"
)

// AuthMiddleware extracts user info from JWT and sets it in context.
// Does NOT block unauthenticated requests - just sets user info if available.
func AuthMiddleware(jwtMgr *jwt.JWTManager) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(AccessTokenCookie)
			if err != nil || cookie.Value == "" {
				// No token, continue as unauthenticated
				return next(c)
			}

			claims, err := jwtMgr.ValidateAccessToken(cookie.Value)
			if err != nil {
				// Token exists but is invalid/expired - signal this to handlers
				// so they can return 401 instead of falling back to anonymous
				c.Set("auth_expired", true)
				return next(c)
			}

			// Set user info in context
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)

			return next(c)
		}
	}
}

// RequireAuth middleware blocks unauthenticated requests.
func RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userID := c.Get("user_id")
			if userID == nil || userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "authentication required",
				})
			}
			return next(c)
		}
	}
}
