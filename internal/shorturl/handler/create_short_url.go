package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

type CreateShortURLResponse struct {
	Message string `json:"message"`
}

func (h *ShortURLHandler) CreateShortURL(c echo.Context) error {
	return c.JSON(http.StatusCreated, CreateShortURLResponse{
		Message: "Short url created successfully",
	})
}
