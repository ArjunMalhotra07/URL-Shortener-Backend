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
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "code is required"})
	}

	output, err := h.Svc.GetLongURL(c.Request().Context(), service.GetLongURLInput{
		Code: code,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCode):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLNotFound):
			return c.JSON(http.StatusNotFound, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLExpired):
			return c.JSON(http.StatusGone, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLInactive):
			return c.JSON(http.StatusGone, ErrorRes{Error: err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	return c.Redirect(http.StatusMovedPermanently, output.LongURL)
}
