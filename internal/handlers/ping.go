package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h *Handlers) PingHandler(c echo.Context) error {
	ctx := c.Request().Context()
	err := h.Usecases.CheckConnection(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
