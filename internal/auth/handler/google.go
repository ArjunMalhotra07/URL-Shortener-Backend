package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/auth/service"
)

type GoogleLoginReq struct {
	IDToken string `json:"id_token"`
}

type GoogleUserInfo struct {
	Sub      string `json:"sub"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified string `json:"email_verified"`
}

func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	var req GoogleLoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.IDToken == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "id_token is required"})
	}

	userInfo, err := verifyGoogleIDToken(req.IDToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "invalid google token"})
	}

	if userInfo.Verified != "true" {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "email not verified"})
	}

	anonID := ""
	if cookie, err := c.Cookie(AnonIDCookie); err == nil && cookie.Value != "" {
		anonID = parseAnonID(cookie.Value)
	}

	output, err := h.Svc.GoogleLogin(c.Request().Context(), service.GoogleLoginInput{
		GoogleID:  userInfo.Sub,
		Email:     userInfo.Email,
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
		AnonID:    anonID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailExistsWithPassword):
			return c.JSON(http.StatusConflict, ErrorRes{Error: "Email already registered with password."})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "google login failed"})
		}
	}

	setAuthCookies(c, output)
	clearAnonCookie(c)

	return c.JSON(http.StatusOK, AuthRes{
		UserID: output.UserID,
		Email:  output.Email,
	})
}

func verifyGoogleIDToken(idToken string) (*GoogleUserInfo, error) {
	resp, err := http.Get("https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid token: " + resp.Status)
	}

	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
