package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "email and password are required"})
	}

	// Get anon_id from cookie for URL transfer
	anonID := ""
	if cookie, err := c.Cookie(AnonIDCookie); err == nil && cookie.Value != "" {
		anonID = parseAnonID(cookie.Value)
	}

	output, err := h.Svc.Login(c.Request().Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
		AnonID:   anonID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "invalid email or password"})
		case errors.Is(err, service.ErrEmailExistsWithGoogle):
			return c.JSON(http.StatusConflict, ErrorRes{Error: "this email uses Google Sign-In. Please login with Google."})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "login failed"})
		}
	}

	setAuthCookies(c, output)
	clearAnonCookie(c)

	return c.JSON(http.StatusOK, AuthRes{
		UserID: output.UserID,
		Email:  output.Email,
	})
}
