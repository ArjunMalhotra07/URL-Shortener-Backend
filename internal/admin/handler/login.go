package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/admin/service"
)

type LoginReq struct {
	AdminID string `json:"admin_id"`
	Code    string `json:"code"`
}

func (h *AdminHandler) Login(c echo.Context) error {
	var req LoginReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.AdminID == "" || req.Code == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "admin_id and code are required"})
	}

	output, err := h.Svc.Login(c.Request().Context(), service.LoginInput{
		AdminID: req.AdminID,
		Code:    req.Code,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "login failed"})
	}

	setAdminCookie(c, output.AccessToken, output.AccessExpiresAt)

	return c.JSON(http.StatusOK, map[string]string{
		"message": "admin login successful",
	})
}
