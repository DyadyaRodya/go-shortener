package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PingHandler godoc
// @Summary      Check server readiness
// @Description  Check server readiness
// @Tags         Debug
// @Param        Cookie header string  false "auth"     default(auth=xxx)
// @Success      200
// @Failure      500
// @Router       /ping [get]
func (h *Handlers) PingHandler(c echo.Context) error {
	ctx := c.Request().Context()
	err := h.Usecases.CheckConnection(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}
