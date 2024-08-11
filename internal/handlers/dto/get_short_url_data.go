package dto

import (
	"encoding/hex"
	"github.com/labstack/echo/v4"
)

const IDParamName = "id"

type GetShortURLData struct {
	ID string
}

func GetShortURLDataFromContext(c echo.Context) (*GetShortURLData, *ErrorResponse) {
	id := c.Param(IDParamName)

	_, err := hex.DecodeString(id)
	if err != nil {
		return nil, ErrInvalidData
	}

	return &GetShortURLData{ID: id}, nil
}
