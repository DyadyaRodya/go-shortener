package handlers

import (
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/url"
)

func (h *Handlers) CreateShortURL(c echo.Context) error {
	createShortURLData, errorResponse := dto.CreateShortURLDataFromContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	ctx := c.Request().Context()
	shortURL, err := h.Usecases.CreateShortURL(ctx, createShortURLData.URL)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusCreated, fullShortURL)
}

func (h *Handlers) CreateShortURLJSON(c echo.Context) error {
	createShortURLData, errorResponse := dto.CreateShortURLDataFromJSONContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	ctx := c.Request().Context()
	shortURL, err := h.Usecases.CreateShortURL(ctx, createShortURLData.URL)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	response := dto.NewCreateShortURLDataResponse(fullShortURL)
	return c.JSON(http.StatusCreated, response)
}
