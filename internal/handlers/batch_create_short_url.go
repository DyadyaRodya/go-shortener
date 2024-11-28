package handlers

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/DyadyaRodya/go-shortener/internal/handlers/dto"
)

// BatchCreateShortURLJSON godoc
// @Summary      Batch create short URLs for full URLs
// @Description  Batch create short URLs for full URLs
// @Tags         Info
// @Accept       json
// @Produce      json
// @Param        Cookie header string  false "auth"     default(auth=xxx)
// @Param        request   body      dto.BatchCreateDataRequest true "Batch create short URLs request"
// @Success      201  {array} dto.BatchCreateDataResponse
// @Failure      400
// @Router       /api/shorten/batch [post]
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
