package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"url_shortner_backend/internal/admin/service"
)

type BlockUserReq struct {
	UserID  string `json:"user_id"`
	Blocked bool   `json:"blocked"`
}

func (h *AdminHandler) BlockUser(c echo.Context) error {
	var req BlockUserReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.UserID == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "user_id is required"})
	}

	err := h.Svc.SetUserBlocked(c.Request().Context(), service.SetUserBlockedInput{
		UserID:  req.UserID,
		Blocked: req.Blocked,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to update user"})
	}

	action := "blocked"
	if !req.Blocked {
		action = "unblocked"
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "user " + action + " successfully",
	})
}
