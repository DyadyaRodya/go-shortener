package dto

import "github.com/DyadyaRodya/go-shortener/internal/domain/entity"

type BatchCreateRequest struct {
	CorrelationID string
	OriginalURL   string
}

type BatchCreateResponse struct {
	CorrelationID string
	ShortURL      *entity.ShortURL
}
