package dto

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type CreateShortURLData struct {
	URL string
}

func CreateShortURLDataFromRequest(r *http.Request) (*CreateShortURLData, *ErrorResponse) {
	if r.Method != http.MethodPost {
		return nil, ErrMethodNotAllowed
	}

	if contentType := r.Header.Get("Content-Type"); !strings.Contains(contentType, "text/plain") {
		return nil, ErrContentType
	}

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

	return &CreateShortURLData{URL: sourceURL}, nil
}
