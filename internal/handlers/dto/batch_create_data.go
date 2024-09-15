package dto

import (
	"encoding/json"
	usecasesdto "github.com/DyadyaRodya/go-shortener/internal/usecases/dto"
	"github.com/labstack/echo/v4"
	"net/url"
	"slices"
)

type BatchCreateDataRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

func BatchCreateDataRequestFromJSONContext(c echo.Context) ([]*BatchCreateDataRequest, *ErrorResponse) {
	data := make([]*BatchCreateDataRequest, 0)

	r := c.Request()
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		return nil, ErrBadData
	}

	correlations := make([]string, 0, len(data))
	for _, e := range data {
		if slices.Contains(correlations, e.CorrelationID) {
			return nil, ErrInvalidData
		}
		correlations = append(correlations, e.CorrelationID)
		if _, err := url.ParseRequestURI(e.OriginalURL); err != nil {
			return nil, ErrInvalidData
		}
	}

	return data, nil
}

func ConvertBatchCreateRequest(dto []*BatchCreateDataRequest) []*usecasesdto.BatchCreateRequest {
	res := make([]*usecasesdto.BatchCreateRequest, 0, len(dto))
	for _, d := range dto {
		res = append(res, &usecasesdto.BatchCreateRequest{
			CorrelationID: d.CorrelationID,
			OriginalURL:   d.OriginalURL,
		})
	}
	return res
}

type BatchCreateDataResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
