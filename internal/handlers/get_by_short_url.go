package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) GetByShortURL(c echo.Context) error {
	getShortURLData, errorResponse := dto.GetShortURLDataFromContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	shortURL, err := h.Usecases.GetShortURL(getShortURLData.ID)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Set("Location", shortURL.URL)
	return c.NoContent(http.StatusTemporaryRedirect)
}
