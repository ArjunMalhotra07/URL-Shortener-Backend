package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

type URLItem struct {
	Code      string `json:"code"`
	LongURL   string `json:"long_url"`
	IsActive  bool   `json:"is_active"`
	CreatedAt string `json:"created_at"`
}

type GetMyURLsRes struct {
	URLs []URLItem `json:"urls"`
}

func (h *ShortURLHandler) GetMyURLs(c echo.Context) error {
	// Get anon_id from cookie
	anonID := getOrCreateAnonID(c)

	// TODO: Check if user is authenticated (JWT), use user_id instead
	ownerType := "anonymous"
	ownerID := anonID

	output, err := h.Svc.GetMyURLs(c.Request().Context(), service.GetMyURLsInput{
		OwnerType: ownerType,
		OwnerID:   ownerID,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidOwner):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLFetch):
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to fetch urls"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	urls := make([]URLItem, len(output.URLs))
	for i, u := range output.URLs {
		urls[i] = URLItem{
			Code:      u.Code,
			LongURL:   u.LongURL,
			IsActive:  u.IsActive,
			CreatedAt: u.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return c.JSON(http.StatusOK, GetMyURLsRes{URLs: urls})
}
