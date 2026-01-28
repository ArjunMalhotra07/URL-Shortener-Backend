package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

func (h *ShortURLHandler) GetOriginalURL(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=invalid")
	}

	output, err := h.Svc.GetLongURL(c.Request().Context(), service.GetLongURLInput{
		Code: code,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCode):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=invalid")
		case errors.Is(err, service.ErrURLNotFound):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=not_found")
		case errors.Is(err, service.ErrURLExpired):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=expired")
		case errors.Is(err, service.ErrURLInactive):
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=inactive")
		default:
			return c.Redirect(http.StatusFound, h.FrontendURL+"/error?type=default")
		}
	}

	return c.Redirect(http.StatusMovedPermanently, output.LongURL)
}
