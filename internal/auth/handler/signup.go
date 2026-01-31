package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

type SignupReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Signup(c echo.Context) error {
	var req SignupReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.Email == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "email and password are required"})
	}

	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "password must be at least 8 characters"})
	}

	anonID := ""
	if cookie, err := c.Cookie(AnonIDCookie); err == nil && cookie.Value != "" {
		anonID = parseAnonID(cookie.Value)
	}

	output, err := h.Svc.Signup(c.Request().Context(), service.SignupInput{
		Email:    req.Email,
		Password: req.Password,
		AnonID:   anonID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExists):
			return c.JSON(http.StatusConflict, ErrorRes{Error: "email already exists"})
		case errors.Is(err, service.ErrEmailExistsWithGoogle):
			return c.JSON(http.StatusConflict, ErrorRes{Error: "this email uses Google Sign-In. Please login with Google."})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to create account"})
		}
	}

	setAuthCookies(c, output)
	clearAnonCookie(c)

	return c.JSON(http.StatusCreated, AuthRes{
		UserID: output.UserID,
		Email:  output.Email,
		Message: "Success Enjoy TinyCLK!",
	})
}
