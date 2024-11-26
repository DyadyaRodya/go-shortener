package dto

import (
	"net/http"
)

// ErrorResponse Custom structure for errors from dto extraction functions
type ErrorResponse struct {
	Message string
	Code    int
}

// Common errors
var (

	// ErrBadData Indicates data parsing errors.
	ErrBadData = &ErrorResponse{Code: http.StatusBadRequest, Message: "BAD_DATA"}

	// ErrInvalidData Indicates data validation errors.
	ErrInvalidData = &ErrorResponse{Code: http.StatusBadRequest, Message: "INVALID_DATA"}
)
