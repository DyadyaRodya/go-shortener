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
	apiGetUserURLs     = "/api/user/urls"
	apiDeleteUserURLs  = "/api/user/urls"
	idURL              = "/:" + dto.IDParamName
	pingURL            = "/ping"
)

func setupRoutes(e *echo.Echo, handlers *handlers.Handlers) {
	e.GET(idURL, handlers.GetByShortURL)
	e.POST(rootURL, handlers.CreateShortURL)
	e.POST(apiShortenURL, handlers.CreateShortURLJSON)
	e.POST(apiShortenURLBatch, handlers.BatchCreateShortURLJSON)
	e.GET(apiGetUserURLs, handlers.GetUserShortURLs)
	e.DELETE(apiDeleteUserURLs, handlers.DeleteUserShortURLs)
	e.GET(pingURL, handlers.PingHandler)
}
