package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) DeleteUserShortURLs(c echo.Context) error {
	ids, errorResponse := dto.IDsFromContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}
	authorized, ok := c.Get("authorized").(bool)
	if !ok || !authorized {
		return c.NoContent(http.StatusUnauthorized)
	}

	userUUID, ok := c.Get("userUUID").(string)
	if !ok || userUUID == "" {
		return c.NoContent(http.StatusUnauthorized)
	}

	req := &usecasesdto.DeleteUserShortURLsRequest{
		UserUUID:      userUUID,
		ShortURLUUIDs: ids,
	}
	go func() { h.DelChan <- req }()

	return c.NoContent(http.StatusAccepted)
}
