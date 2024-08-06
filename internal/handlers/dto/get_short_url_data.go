package dto

import (
	"encoding/hex"
	"net/http"
	"strings"
)

type GetShortURLData struct {
	ID string
}

func GetShortURLDataFromRequest(r *http.Request) (*GetShortURLData, *ErrorResponse) {
	if r.Method != http.MethodGet {
		return nil, ErrMethodNotAllowed
	}

	id := strings.TrimPrefix(r.URL.Path, "/")

	_, err := hex.DecodeString(id)
	if err != nil {
		return nil, ErrInvalidData
	}

	return &GetShortURLData{ID: id}, nil
}
