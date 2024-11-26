package handlers

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/domain/entity"
	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
)

func (h *Handlers) CreateShortURL(c echo.Context) error {
	createShortURLData, errorResponse := dto.CreateShortURLDataFromContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	userUUID, ok := c.Get("userUUID").(string)
	if !ok {
		userUUID = ""
	}

	ctx := c.Request().Context()
	shortURL, err := h.Usecases.CreateShortURL(ctx, createShortURLData.URL, userUUID)
	if err != nil && !errors.Is(err, entity.ErrShortURLExists) {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var statusCode int
	if err != nil {
		statusCode = http.StatusConflict
	} else {
		statusCode = http.StatusCreated
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(statusCode, fullShortURL)
}

func (h *Handlers) CreateShortURLJSON(c echo.Context) error {
	createShortURLData, errorResponse := dto.CreateShortURLDataFromJSONContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	userUUID, ok := c.Get("userUUID").(string)
	if !ok {
		userUUID = ""
	}

	ctx := c.Request().Context()
	shortURL, err := h.Usecases.CreateShortURL(ctx, createShortURLData.URL, userUUID)
	if err != nil && !errors.Is(err, entity.ErrShortURLExists) {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	var statusCode int
	if err != nil {
		statusCode = http.StatusConflict
	} else {
		statusCode = http.StatusCreated
	}

	fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, shortURL.ID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	response := dto.NewCreateShortURLDataResponse(fullShortURL)
	return c.JSON(statusCode, response)
}
