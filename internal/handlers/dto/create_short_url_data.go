package dto

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

type CreateShortURLDataRequest struct {
	URL string `json:"url"`
}

func CreateShortURLDataFromContext(c echo.Context) (*CreateShortURLDataRequest, *ErrorResponse) {
	r := c.Request()

	defer func() {
		_ = r.Body.Close()
	}()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, ErrBadData
	}

	sourceURL := strings.TrimSpace(string(body))
	_, err = url.ParseRequestURI(sourceURL)
	if err != nil {
		return nil, ErrInvalidData
	}

	return &CreateShortURLDataRequest{URL: sourceURL}, nil
}

func CreateShortURLDataFromJSONContext(c echo.Context) (*CreateShortURLDataRequest, *ErrorResponse) {
	data := &CreateShortURLDataRequest{}

	r := c.Request()
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return nil, ErrBadData
	}

	if _, err := url.ParseRequestURI(data.URL); err != nil {
		return nil, ErrInvalidData
	}

	return data, nil
}

type CreateShortURLDataResponse struct {
	Result string `json:"result"`
}

func NewCreateShortURLDataResponse(URL string) *CreateShortURLDataResponse {
	return &CreateShortURLDataResponse{Result: URL}
}
