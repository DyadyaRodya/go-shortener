package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
)

// GetStats godoc
// @Summary      Get total numbers  of users and shortened URLs
// @Description  Get total numbers  of users and shortened URLs
// @Tags         Info Internal
// @Produce      json
// @Param        X-Real-IP header string true "real ip address from webserver" default()
// @Success      200  {object} dto.StatsData
// @Failure      403
// @Router       /api/internal/stats [get]
func (h *Handlers) GetStats(c echo.Context) error {
	ctx := c.Request().Context()
	stats, err := h.Usecases.GetStats(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	resp := dto.NewStatsDataResponse(stats)
	return c.JSON(http.StatusOK, resp)
}
