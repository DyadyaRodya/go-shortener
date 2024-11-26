package dto

import (
	"encoding/hex"

	"github.com/labstack/echo/v4"
)

// IDParamName Name of path parameter with short URL ID
const IDParamName = "id"

// GetShortURLData Structure for querying short URL by ID
type GetShortURLData struct {
	ID string
}

// GetShortURLDataFromContext Extracts short URL query from request
func GetShortURLDataFromContext(c echo.Context) (*GetShortURLData, *ErrorResponse) {
	id := c.Param(IDParamName)

	_, err := hex.DecodeString(id)
	if err != nil {
		return nil, ErrInvalidData
	}

	return &GetShortURLData{ID: id}, nil
}
