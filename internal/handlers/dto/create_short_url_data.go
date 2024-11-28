package dto

import (
	"encoding/json"
	"io"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
)

// CreateShortURLDataRequest Structure of JSON body for create shortURL requests. URL is fullURL
type CreateShortURLDataRequest struct {
	URL string `json:"url"`
}

// CreateShortURLDataFromContext Extracts create short URL request data from plain text body
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

// CreateShortURLDataFromJSONContext Extracts create short URL request data from JSON body
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

// CreateShortURLDataResponse Structure of JSON body for create shortURL response. Result is complete short URL
type CreateShortURLDataResponse struct {
	Result string `json:"result"`
}

// NewCreateShortURLDataResponse Creates *CreateShortURLDataResponse with given short URL
func NewCreateShortURLDataResponse(URL string) *CreateShortURLDataResponse {
	return &CreateShortURLDataResponse{Result: URL}
}
