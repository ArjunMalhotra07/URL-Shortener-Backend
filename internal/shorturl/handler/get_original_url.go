package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetOriginalURLRes struct {
	Message string `json:"message"`
}

func (h *ShortURLHandler) GetOriginalURL(c echo.Context) error {
	return c.JSON(http.StatusMovedPermanently, GetOriginalURLRes{
		Message: "Short url created successfully",
	})
}
