package app

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
)

const (
	rootURL = "/"
	idURL   = "/:" + dto.IDParamName
)

func setupRoutes(e *echo.Echo, handlers *handlers.Handlers) {
	e.GET(idURL, handlers.GetByShortURL)
	e.POST(rootURL, handlers.CreateShortURL)
}
