package app

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
)

const (
	rootURL            = "/"
	apiShortenURL      = "/api/shorten"
	apiShortenURLBatch = "/api/shorten/batch"
	idURL              = "/:" + dto.IDParamName
	pingURL            = "/ping"
)

func setupRoutes(e *echo.Echo, handlers *handlers.Handlers) {
	e.GET(idURL, handlers.GetByShortURL)
	e.POST(rootURL, handlers.CreateShortURL)
	e.POST(apiShortenURL, handlers.CreateShortURLJSON)
	e.POST(apiShortenURLBatch, handlers.BatchCreateShortURLJSON)
	e.GET(pingURL, handlers.PingHandler)
}
