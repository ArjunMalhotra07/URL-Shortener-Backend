package handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type MeRes struct {
	UserID               string  `json:"user_id"`
	Email                string  `json:"email"`
	Name                 string  `json:"name,omitempty"`
	AvatarURL            string  `json:"avatar_url,omitempty"`
	Tier                 string  `json:"tier"`
	SubscriptionEndsAt   *string `json:"subscription_ends_at,omitempty"`
	URLsCreatedThisMonth int64   `json:"urls_created_this_month"`
	URLsLimit            int     `json:"urls_limit"`
}

func (h *AuthHandler) Me(c echo.Context) error {
	userID := c.Get("user_id")

	if userID == nil || userID == "" {
		return c.JSON(http.StatusUnauthorized, ErrorRes{Error: "not authenticated"})
	}

	output, err := h.Svc.GetMe(c.Request().Context(), userID.(string))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to get user info"})
	}

	res := MeRes{
		UserID:               output.UserID,
		Email:                output.Email,
		Name:                 output.Name,
		AvatarURL:            output.AvatarURL,
		Tier:                 output.Tier,
		URLsCreatedThisMonth: output.URLsCreatedThisMonth,
		URLsLimit:            output.URLsLimit,
	}

	if output.SubscriptionEndsAt != nil {
		formatted := output.SubscriptionEndsAt.UTC().Format(time.RFC3339)
		res.SubscriptionEndsAt = &formatted
	}

	return c.JSON(http.StatusOK, res)
}
