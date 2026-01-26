package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type CreateShortURLRes struct {
	Message string `json:"message"`
}

func (h *ShortURLHandler) CreateShortURL(c echo.Context) error {
	return c.JSON(http.StatusCreated, CreateShortURLRes{
		Message: "Short url created successfully",
	})
}
