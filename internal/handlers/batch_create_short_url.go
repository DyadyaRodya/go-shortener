package handlers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
)

func (h *Handlers) BatchCreateShortURLJSON(c echo.Context) error {
	batchCreateShortURLData, errorResponse := dto.BatchCreateDataRequestFromJSONContext(c)
	if errorResponse != nil {
		return c.NoContent(errorResponse.Code)
	}

	ctx := c.Request().Context()
	req := dto.ConvertBatchCreateRequest(batchCreateShortURLData)

	userUUID, ok := c.Get("userUUID").(string)
	if !ok {
		userUUID = ""
	}

	resp, err := h.Usecases.BatchCreateShortURLs(ctx, req, userUUID)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	result := make([]*dto.BatchCreateDataResponse, 0, len(resp))
	for _, d := range resp {
		fullShortURL, err := url.JoinPath(h.Config.BaseShortURL, d.ShortURL.ID)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		result = append(result, &dto.BatchCreateDataResponse{
			CorrelationID: d.CorrelationID,
			ShortURL:      fullShortURL,
		})
	}
	return c.JSON(http.StatusCreated, result)
}
