package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
)

// GetByShortURL godoc
// @Summary      Get full URL for short URL
// @Description  Get full URL for short URL
// @Tags         Info
// @Param        Cookie header string  false "auth"     default(auth=xxx)
// @Param        id   path     string  true "short URL ID"
// @Success      307
// @Header       307  {string} Location "Full URL"
// @Failure      400
// @Failure      410
// @Router       /{id} [get]
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
