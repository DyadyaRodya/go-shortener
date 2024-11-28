package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
)

// DeleteUserShortURLs godoc
// @Summary      Delete user short URLs
// @Description  Delete user short URLs
// @Tags         Info
// @Accept       json
// @Param        Cookie header string  false "auth"     default(auth=xxx)
// @Param        request   body      []string true "Delete user short URLs request"
// @Success      202
// @Failure      400
// @Failure      401
// @Router       /api/user/urls [delete]
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
