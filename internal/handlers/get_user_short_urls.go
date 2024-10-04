package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

func (h *Handlers) GetUserShortURLs(c echo.Context) error {
	authorized, ok := c.Get("authorized").(bool)
	if !ok || !authorized {
		return c.NoContent(http.StatusUnauthorized)
	}

	userUUID, ok := c.Get("userUUID").(string)
	if !ok || userUUID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	ctx := c.Request().Context()
	shortURLs, err := h.Usecases.GetUserShortURLs(ctx, userUUID)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	if len(shortURLs) == 0 {
		return c.NoContent(http.StatusNoContent)
	}

	resp := make([]*dto.ShortURLData, 0, len(shortURLs))
	for _, shortURL := range shortURLs {
		fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		resp = append(resp, &dto.ShortURLData{ShortURL: fullShortURL, OriginalURL: shortURL.URL})
	}

	return c.JSON(http.StatusOK, resp)
}
