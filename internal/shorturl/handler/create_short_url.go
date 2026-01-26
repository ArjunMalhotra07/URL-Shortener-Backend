package handler

import (
	"errors"
	"net/http"

	"url_shortner_backend/internal/shorturl/service"

	"github.com/labstack/echo/v4"
)

type CreateShortURLReq struct {
	LongURL string `json:"long_url"`
}

type CreateShortURLRes struct {
	Code    string `json:"code"`
	LongURL string `json:"long_url"`
}

type ErrorRes struct {
	Error string `json:"error"`
}

func (h *ShortURLHandler) CreateShortURL(c echo.Context) error {
	var req CreateShortURLReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "invalid request body"})
	}

	if req.LongURL == "" {
		return c.JSON(http.StatusBadRequest, ErrorRes{Error: "long_url is required"})
	}

	output, err := h.Svc.CreateShortURL(c.Request().Context(), service.CreateShortURLInput{
		LongURL: req.LongURL,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			return c.JSON(http.StatusBadRequest, ErrorRes{Error: err.Error()})
		case errors.Is(err, service.ErrURLCreation), errors.Is(err, service.ErrURLCodeUpdate):
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "failed to create short url"})
		default:
			return c.JSON(http.StatusInternalServerError, ErrorRes{Error: "internal server error"})
		}
	}

	return c.JSON(http.StatusCreated, CreateShortURLRes{
		Code:    output.Code,
		LongURL: output.LongURL,
	})
}
