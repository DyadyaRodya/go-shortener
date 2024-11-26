package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handlers) PingHandler(c echo.Context) error {
	ctx := c.Request().Context()
	err := h.Usecases.CheckConnection(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
