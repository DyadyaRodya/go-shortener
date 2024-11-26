package dto

import (
	"encoding/hex"
	"encoding/json"

	"github.com/labstack/echo/v4"
)

func IDsFromContext(c echo.Context) ([]string, *ErrorResponse) {
	var ids []string

	r := c.Request()
	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		return nil, ErrBadData
	}

	for _, id := range ids {
		_, err := hex.DecodeString(id)
		if err != nil {
			return nil, ErrInvalidData
		}
	}
	return ids, nil
}
