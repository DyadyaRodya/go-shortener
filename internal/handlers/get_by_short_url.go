package handlers

import (
	"errors"
	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) GetByShortURL(c echo.Context) error {
	getShortURLData, errorResponse := dto.GetShortURLDataFromContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	ctx := c.Request().Context()
	shortURL, err := h.Usecases.GetShortURL(ctx, getShortURLData.ID)
	if err != nil {
		if errors.Is(err, entity.ErrShortURLDeleted) {
			return c.NoContent(http.StatusGone)
		}
		return c.NoContent(http.StatusBadRequest)
	}

	c.Response().Header().Set("Location", shortURL.URL)
	return c.NoContent(http.StatusTemporaryRedirect)
}
