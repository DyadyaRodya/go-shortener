package dto

import (
	"net/http"
)

type ErrorResponse struct {
	Message string
	Code    int
}

var (
	ErrMethodNotAllowed = &ErrorResponse{Code: http.StatusMethodNotAllowed, Message: "METHOD_NOT_ALLOWED"}
	ErrContentType      = &ErrorResponse{Code: http.StatusUnsupportedMediaType, Message: "UNSUPPORTED_CONTENT_TYPE"}
	ErrBadData          = &ErrorResponse{Code: http.StatusBadRequest, Message: "BAD_DATA"}
	ErrInvalidData      = &ErrorResponse{Code: http.StatusBadRequest, Message: "INVALID_DATA"}
)
